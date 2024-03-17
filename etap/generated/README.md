# Generated eTapestry Client

This directory contains the eTapestry client, extracted from the SOAP API into a Go client.

## Regenerating the Code

Step 1: download the WSDL definition: 

```bash
curl https://sna.etapestry.com/v3messaging/service?WSDL -o etapestry-api.wsdl
```

Step 2: Install the [wsdl2go](https://github.com/fiorix/wsdl2go) tool, either by downloading a [pre-built binary from the release](https://github.com/fiorix/wsdl2go/releases), or with:

```bash
go install github.com/fiorix/wsdl2go@latest
```

Step 3: Use `wsdl2go` to regenerate this package:

```
wsdl2go -p generated < etapestry-api.wsdl > generated.go
```

## Known Bugs

**NOTE** The generated code is _almost_ correct out the box, but has a few compile time errors in Go 1.22:

- Replace `&DateTime{}` ==> `""` in a single location.
- Add **ArrayType xml.Attr `xml:"arrayType,attr"`** to `arrayofint`

Once the code compiles, you're good to go.