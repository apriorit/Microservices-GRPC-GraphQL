package server

import (
	holdersv1 "tutorial/gen/go/proto/holders"
)

var (
	testholder1 = &holdersv1.Holder{
		FirstName: "Silvia",
		LastName:  "McClintock",
		Phone:     "916-335-5759",
		Email:     "SilviaAMcClintock@dayrep.com",
	}
	testholder2 = &holdersv1.Holder{
		FirstName: "Carl",
		LastName:  "Koch",
		Phone:     "515-276-5355",
		Email:     "CarlRKoch@jourrapide.com",
	}
)
