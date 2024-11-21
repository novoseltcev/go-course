package main

func main() {
	if err := Cmd().Execute(); err != nil {
		panic(err)
	}
}
