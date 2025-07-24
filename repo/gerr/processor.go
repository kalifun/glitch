package gerr

import "context"

// Middleware allows processing errors in a chain
type Middleware func(ctx context.Context, err *Error, next func(context.Context, *Error) interface{}) interface{}

// ProcessorEngine is the main error processing engine
type ProcessorEngine struct {
	registry   Registry
	localizer  Localizer
	formatter  Formatter
	handlers   []Handler
	middleware []Middleware
}

// NewProcessorEngine creates a new error processing engine
func NewProcessorEngine() *ProcessorEngine {
	return &ProcessorEngine{
		registry:   globalRegistry,
		localizer:  NewDefaultLocalizer(),
		formatter:  NewDefaultFormatter(),
		handlers:   make([]Handler, 0),
		middleware: make([]Middleware, 0),
	}
}

// SetRegistry sets the error registry
func (e *ProcessorEngine) SetRegistry(registry Registry) *ProcessorEngine {
	e.registry = registry
	return e
}

// SetLocalizer sets the localizer
func (e *ProcessorEngine) SetLocalizer(localizer Localizer) *ProcessorEngine {
	e.localizer = localizer
	return e
}

// SetFormatter sets the formatter
func (e *ProcessorEngine) SetFormatter(formatter Formatter) *ProcessorEngine {
	e.formatter = formatter
	return e
}

// AddHandler adds an error handler
func (e *ProcessorEngine) AddHandler(handler Handler) *ProcessorEngine {
	e.handlers = append(e.handlers, handler)
	return e
}

// AddMiddleware adds middleware to the processing chain
func (e *ProcessorEngine) AddMiddleware(middleware Middleware) *ProcessorEngine {
	e.middleware = append(e.middleware, middleware)
	return e
}

// Process processes an error through the engine
func (e *ProcessorEngine) Process(ctx context.Context, err *Error) interface{} {
	// Apply middleware chain
	var next func(context.Context, *Error) interface{}
	next = func(ctx context.Context, err *Error) interface{} {
		// Try handlers first
		for _, handler := range e.handlers {
			if handler.CanHandle(err) {
				return handler.Handle(ctx, err)
			}
		}

		// Default processing
		return e.formatter.Format(ctx, err)
	}

	// Apply middleware in reverse order
	for i := len(e.middleware) - 1; i >= 0; i-- {
		middleware := e.middleware[i]
		prevNext := next
		next = func(ctx context.Context, err *Error) interface{} {
			return middleware(ctx, err, prevNext)
		}
	}

	return next(ctx, err)
}

// Global engine instance
var globalEngine *ProcessorEngine

func init() {
	globalEngine = NewProcessorEngine()
	// Set up the localizer with the registry
	if localizer, ok := globalEngine.localizer.(*DefaultLocalizer); ok {
		localizer.SetRegistry(globalEngine.registry)
	}
}

// Process processes an error using the global engine
func Process(ctx context.Context, err *Error) interface{} {
	return globalEngine.Process(ctx, err)
}
