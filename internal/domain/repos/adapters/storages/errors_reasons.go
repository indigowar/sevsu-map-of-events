package storages

/**

This file contains a set of error-reasons.
This integer values are used to change logic of handling error.
When the repository returns an error, you can access to this reason as:
```

_, errors := repository.Job()

switch (errors.Reason) {

case ErrReason1:
	{ do here a specific logic }
	break
case ErrReason2:
	{ do here a specific logic }

and so on.
}

```

 **/

const (
	ErrReasonInternalStorageErr = iota
	ErrReasonObjectNotFoundErr
	ErrReasonObjectAlreadyExistsErr
)
