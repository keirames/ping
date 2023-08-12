package customerror

import "github.com/vektah/gqlparser/v2/gqlerror"

func BadRequest() *gqlerror.Error {
	return gqlerror.Errorf("Bad Request")
}

func NotFound() *gqlerror.Error {
	return gqlerror.Errorf("Not Found")
}
