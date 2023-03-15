package internal

type (
	Filter interface {
		CheckOnlyOfficial()
		CheckFlags()
		CheckOnlyUserId()
	}

	BaseFilter struct {
		MustFlags    *PostFlag
		MustNotFlags *PostFlag
		*FilterOptions
	}

	FilterOptions struct {
		OnlyUserId   *string
		OnlyOfficial *bool
	}
)

func (f *BaseFilter) CheckOnlyOfficial() {
	if f.OnlyOfficial != nil {
		f.MustFlags = f.MustFlags.SetFlag(OfficialFlag, *f.OnlyOfficial)
	}
}
