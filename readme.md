## Twikle

Twikle is a small utility that applies a pattern to png images.

![demo](demo/demo.png)

> If your image is not a png, you can use the included script `toPng.sh` to convert it

#### Build

```bash
docker build . -t twikle
```

#### Usage

```bash
./twikle.sh -i Lenna.png --pattern=brick.png --patternscale=20
```

Input images must be in the folder `images/input/`, patterns must be in the folder `images/pattern/`. Output images will be written in `images/output/`. 

Help is available with `./twikle.sh -h`.

#### Dev

```bash
docker build . -f dockerfile_dev -t twikle_dev
docker run -it --rm -v "$(pwd)/src":"/go/src" -v "$(pwd)/images":"/images" twikle_dev
```

Then, in the container :

```bash
cd src/twikle
go run twikle.go -i vault_boy.png --pattern=brick.png --patternscale=20
```
