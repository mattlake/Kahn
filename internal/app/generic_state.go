package app

// GenericConfirmationState provides a generic confirmation dialog state management
type GenericConfirmationState[T string] struct {
	showConfirm  bool
	itemToDelete T
	errorMessage string
}

// NewGenericConfirmationState creates a new generic confirmation state
func NewGenericConfirmationState[T string]() *GenericConfirmationState[T] {
	return &GenericConfirmationState[T]{}
}

// ShowConfirm shows confirmation dialog for the specified item
func (cs *GenericConfirmationState[T]) ShowConfirm(item T) {
	cs.showConfirm = true
	cs.itemToDelete = item
	cs.errorMessage = ""
}

// HideConfirm hides confirmation dialog and clears state
func (cs *GenericConfirmationState[T]) HideConfirm() {
	cs.showConfirm = false
	cs.itemToDelete = ""
	cs.errorMessage = ""
}

// IsShowingConfirm returns true if confirmation dialog is showing
func (cs *GenericConfirmationState[T]) IsShowingConfirm() bool {
	return cs.showConfirm
}

// GetItemToDelete returns the item awaiting confirmation
func (cs *GenericConfirmationState[T]) GetItemToDelete() T {
	return cs.itemToDelete
}

// SetError sets an error message
func (cs *GenericConfirmationState[T]) SetError(message string) {
	cs.errorMessage = message
}

// GetError returns the error message
func (cs *GenericConfirmationState[T]) GetError() string {
	return cs.errorMessage
}

// HasError returns true if there is an error
func (cs *GenericConfirmationState[T]) HasError() bool {
	return cs.errorMessage != ""
}

// ClearError clears the error message
func (cs *GenericConfirmationState[T]) ClearError() {
	cs.errorMessage = ""
}
