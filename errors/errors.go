package errors

import (
	"fmt"
	"math/rand"
	"net/http"
)

func GlamrErrorGeneralServerError(reason string) GlamrError {
	return glamrError{
		reason:     reason,
		httpStatus: http.StatusInternalServerError,
	}
}

func GlamrErrorGeneralNotFound(reason string) GlamrError {
	return glamrError{
		reason:     reason,
		httpStatus: http.StatusNotFound,
	}
}

func GlamrErrorGeneralBadRequest(reason string) GlamrError {
	return glamrError{
		reason:     reason,
		httpStatus: http.StatusBadRequest,
	}
}

func GlamrUnauthenticated() GlamrError {
	reasons := []string{
		"Unauthenticated",
	}
	return glamrError{
		reason:     reasons[rand.Int()%len(reasons)],
		httpStatus: http.StatusUnauthorized,
	}
}

func GlamrDeactivated() GlamrError {
	reasons := []string{
		"The account has been deleted.",
	}
	return glamrError{
		reason:     reasons[rand.Int()%len(reasons)],
		httpStatus: http.StatusUnauthorized,
	}
}

func GlamrErrorBadRequest() GlamrError {
	reasons := []string{
		"Oops! this I don't understand at all.",
		"Your request is like a pineapple on pizza - just plain wrong.",
		"We're not sure what you were trying to do.",
		"I'm sorry, but I can't hear you over the sound of your malformed request.",
		"Your request is causing the parser to have an existential crisis.",
	}
	return glamrError{
		reason:     reasons[rand.Int()%len(reasons)],
		httpStatus: http.StatusBadRequest,
	}
}

func GlamrErrorInvalidValue(name string) GlamrError {
	reasons := []string{
		"Houston, we have a problem. Invalid value detected in '%s'!",
		"We are a bit confused... '%s' seems invalid.",
		"We're scratching our heads over here... '%s' seems invalid.",
	}
	reason := reasons[rand.Int()%len(reasons)]
	return glamrError{
		reason:     fmt.Sprintf(reason, name),
		httpStatus: http.StatusBadRequest,
	}
}

func GlamrErrorMissingField(name string) GlamrError {
	reasons := []string{
		"Whoops! You forgot something important in your request... '%s'!",
		"We are a bit confused... looks like you missed a field... '%s'.",
		"Houston, we have a problem... it seems like you didn't provide a parameter ... '%s'.",
		"Please provide the missing piece to this puzzle in your request... '%s'!",
		"We're scratching our heads over here... seems like you missed '%s'.",
	}
	reason := reasons[rand.Int()%len(reasons)]
	return glamrError{
		reason:     fmt.Sprintf(reason, name),
		httpStatus: http.StatusBadRequest,
	}
}

func GlamrErrorInternalServerError() GlamrError {
	reasons := []string{
		"Oops. Something went wrong :-(",
	}
	reason := reasons[rand.Int()%len(reasons)]
	return glamrError{
		reason:     reason,
		httpStatus: http.StatusInternalServerError,
	}
}

func GlamrErrorDatabaseIssue() GlamrError {
	reasons := []string{
		"Oops. Something went wrong :-(",
	}
	reason := reasons[rand.Int()%len(reasons)]
	return glamrError{
		reason:     reason,
		httpStatus: http.StatusBadRequest,
	}
}
