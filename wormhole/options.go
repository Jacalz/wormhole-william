package wormhole

type SendOption interface {
	setOption(*sendOptions) error
}

// WithCode returns a SendOption to use a specific nameplate+code
// instead of generating one dynamically.
func WithCode(code string) SendOption {
	return sendCodeOption{code: code}
}

// WithProgress returns a SendOption to track the progress of the data
// transfer. It takes a callback function that will be called for each
// chunk of data successfully written.
//
// WithProgress is only minimally supported in SendText. SendText does
// not use the wormhole transit protocol so it is not able to detect
// the progress of the receiver. This limitation does not apply to
// SendFile or SendDirectory.
func WithProgress(f func(sentBytes int64, totalBytes int64)) SendOption {
	return progressSendOption{f}
}

type sendOptions struct {
	code         string
	progressFunc progressFunc
}

type progressFunc func(sentBytes int64, totalBytes int64)

type progressSendOption struct {
	progressFunc progressFunc
}

func (o progressSendOption) setOption(opts *sendOptions) error {
	opts.progressFunc = o.progressFunc
	return nil
}

type sendCodeOption struct {
	code string
}

func (o sendCodeOption) setOption(opts *sendOptions) error {
	if err := validateCode(o.code); err != nil {
		return err
	}

	opts.code = o.code
	return nil
}
