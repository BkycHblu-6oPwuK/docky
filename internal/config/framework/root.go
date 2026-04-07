package framework

type Framework string

const (
	Bitrix     Framework = "bitrix"
	BitrixNuxt Framework = "bitrix_nuxt"
	Laravel    Framework = "laravel"
	Symfony    Framework = "symfony"
	Yii2	   Framework = "yii2"
	Yii3	   Framework = "yii3"
	Wordpress  Framework = "wordpress"
	Vanilla    Framework = "vanilla"
)

func GetAll() []Framework {
	return []Framework{
		Bitrix,
		BitrixNuxt,
		Laravel,
		Symfony,
		Yii2,
		Yii3,
		Wordpress,
		Vanilla,
	}
}

func GetAllStrings() []string {
	all := GetAll()
	strs := make([]string, len(all))
	for i, f := range all {
		strs[i] = f.String()
	}
	return strs
}

func ParseFramework(s string) Framework {
	switch s {
	case Bitrix.String():
		return Bitrix
	case BitrixNuxt.String():
		return BitrixNuxt
	case Laravel.String():
		return Laravel
	case Symfony.String():
		return Symfony
	case Yii2.String():
		return Yii2
	case Yii3.String():
		return Yii3
	case Wordpress.String():
		return Wordpress
	case Vanilla.String():
		return Vanilla
	default:
		panic("Неизвестный фреймворк")
	}
}

func (f Framework) String() string {
	return string(f)
}
