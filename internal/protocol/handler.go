package protocol

type Response struct {
	Status int
	Retry  bool
	Class  FailureClass
}

type Handler struct{ repository *Repository }

func NewHandler(repository *Repository) *Handler { return &Handler{repository: repository} }

func (h *Handler) Ingest(data []byte) Response {
	frame, err := Decode(data)
	if err == nil {
		err = h.repository.Store(frame)
	}
	if err == nil {
		return Response{Status: 202}
	}
	class := Classify(err)
	switch class {
	case FailurePermanent:
		return Response{Status: 422, Class: class}
	case FailureMissing:
		return Response{Status: 404, Class: class}
	default:
		return Response{Status: 503, Retry: true, Class: class}
	}
}
