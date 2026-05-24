package filevariant

var (
	ErrUnknownVariantType          = NewFileVariantError("domain.file_variant.error.unknown_variant_type")
	ErrMarshalMeta                 = NewFileVariantError("domain.file_variant.error.marshal_meta_error")
	ErrUmarshalMeta                = NewFileVariantError("domain.file_variant.error.unmarshal_meta_error")
	ErrInvalidSize                 = NewFileVariantError("domain.file_variant.error.invalid_size")
	ErrInvalidMimeType             = NewFileVariantError("domain.file_variant.error.invalid_mime_type")
	ErrInvalidMeta                 = NewFileVariantError("domain.file_variant.error.invalid_meta")
	ErrVariantNotAllowedForPurpose = NewFileVariantError("domain.file_variant.error.variant_not_allowed_for_purpose")
)
