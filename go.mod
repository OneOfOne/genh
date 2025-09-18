module go.oneofone.dev/genh

go 1.25

require github.com/vmihailenco/msgpack/v5 v5.3.5

require github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect

replace github.com/vmihailenco/msgpack/v5 v5.3.5 => github.com/alpineiq/msgpack/v5 v5.3.5-no-partial-alloc
