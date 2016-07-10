package main

type Renderer interface {
	Render(code int32, fileName string)
}
