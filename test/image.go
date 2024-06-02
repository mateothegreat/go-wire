package test

type Image struct {
	Camera string `msgpack:"camera"`
	Data   []byte `msgpack:"data"`
}
