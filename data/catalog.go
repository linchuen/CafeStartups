package data

import _ "embed"

// MVPFixture is the provisional 84-card catalog embedded into the server binary.
//
//go:embed mvp-fixture.json
var MVPFixture []byte
