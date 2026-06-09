package service

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"time"
	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/repository"
)

type PalmService struct {
	palmRepo     *repository.PalmRepository
	pairingRepo  *repository.PairingRepository
	cryptoSvc    *CryptoService
}

func NewPalmService(
	palmRepo *repository.PalmRepository,
	pairingRepo *repository.PairingRepository,
	cryptoSvc *CryptoService,
) *PalmService {
	return &PalmService{
		palmRepo:    palmRepo,
		pairingRepo: pairingRepo,
		cryptoSvc:   cryptoSvc,
	}
}

// Math Utilities
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dot float64
	var normA float64
	var normB float64

	for i := range a {
		x := float64(a[i])
		y := float64(b[i])

		dot += x * y
		normA += x * x
		normB += y * y
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func AverageEmbeddings(embeddings [][]float32) []float32 {
	if len(embeddings) == 0 {
		return nil
	}

	dim := len(embeddings[0])
	avg := make([]float32, dim)

	for _, emb := range embeddings {
		for i := 0; i < dim; i++ {
			avg[i] += emb[i]
		}
	}

	for i := 0; i < dim; i++ {
		avg[i] /= float32(len(embeddings))
	}

	return NormalizeEmbedding(avg)
}

func NormalizeEmbedding(embedding []float32) []float32 {
	var norm float64

	for _, v := range embedding {
		norm += float64(v * v)
	}

	norm = math.Sqrt(norm)

	if norm == 0 {
		return embedding
	}

	for i := range embedding {
		embedding[i] = float32(float64(embedding[i]) / norm)
	}

	return embedding
}

func Float32SliceToBytes(f []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, f)
	return buf.Bytes(), err
}

func BytesToFloat32Slice(b []byte) ([]float32, error) {
	if len(b)%4 != 0 {
		return nil, errors.New("invalid byte slice length for float32")
	}
	f := make([]float32, len(b)/4)
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &f)
	return f, err
}

// Enrollment Logic
type EnrollInput struct {
	SessionToken   string
	DeviceCode     string
	HandSide       string
	ModelVersion   string
	EmbeddingDim   int
	Embeddings     [][]float32
	LivenessPassed bool
	QualityScore   float64
	ThermalMin     float64
	ThermalMax     float64
	ThermalAvg     float64
}

func (s *PalmService) EnrollPalm(input EnrollInput) (*model.PalmTemplate, error) {
	// 1. Validate session
	session, err := s.pairingRepo.FindByToken(input.SessionToken)
	if err != nil {
		return nil, errors.New("invalid session token")
	}

	if session.Status != "approved" {
		return nil, errors.New("session is not approved")
	}

	if session.UserID == nil {
		return nil, errors.New("no user linked to session")
	}

	if !input.LivenessPassed {
		return nil, errors.New("liveness check failed")
	}

	// 2. Average embeddings
	avgEmbedding := AverageEmbeddings(input.Embeddings)
	embBytes, err := Float32SliceToBytes(avgEmbedding)
	if err != nil {
		return nil, errors.New("failed to process embeddings")
	}

	// 3. Encrypt
	enc, nonce, err := s.cryptoSvc.Encrypt(embBytes)
	if err != nil {
		return nil, errors.New("failed to encrypt template")
	}

	// Use HandSide from session if available, otherwise fallback to input
	handSide := input.HandSide
	if session.HandSide != nil && *session.HandSide != "" {
		handSide = *session.HandSide
	}

	// 4. Save template
	template := &model.PalmTemplate{
		UserID:             *session.UserID,
		HandSide:           handSide,
		TemplateEncrypted:  enc,
		TemplateNonce:      nonce,
		EmbeddingDim:       input.EmbeddingDim,
		ModelVersion:       input.ModelVersion,
		Status:             "active",
		Threshold:          0.8200,
		RegisteredDeviceID: &session.DeviceID,
	}

	if err := s.palmRepo.Create(template); err != nil {
		return nil, errors.New("failed to save palm template")
	}

	// 5. Complete session
	now := time.Now()
	session.Status = "completed"
	session.CompletedAt = &now
	_ = s.pairingRepo.Update(session)

	return template, nil
}