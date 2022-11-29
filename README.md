# pixgo

pixgo is a library for interacting with Pixoo devices.

This library is heavily based on SomethingWithComputers' [pixoo library](https://github.com/SomethingWithComputers/pixoo).

## Usage

```go
import "github.com/SethCurry/pixgo"

client := pixgo.NewClient("192.168.1.12", 64)

# draw a letter on the screen
err := client.DrawCharacter('a', 0, 0, 255, 255, 255)
if err != nil {
    panic(err.Error())
}

# send the updated buffer to the Pixoo
err = client.Push()
if err != nil {
    panic(err.Error())
}
```