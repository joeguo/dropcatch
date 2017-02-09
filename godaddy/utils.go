package godaddy

import (
	"errors"
)

func DomainId(operation, tld string, period int) (int, error) {
	//can store in database
	if period < 1 || period > 10 {
		return 0, errors.New("Period must between 1 and 10")
	}
	switch tld{
	case "com":
		switch operation{
		case Register:
			return 350001 + (period - 1), nil
		case Renew:
			return 350012 + (period - 1), nil
		case Transfer:
			return 350011, nil
		}
	case "net":
		switch operation{
		case Register:
			return 350030 + (period - 1), nil
		case Renew:
			return 350041 + (period - 1), nil
		case Transfer:
			return 350040, nil
		}
	case "org":
		switch operation{
		case Register:
			return 350150 + (period - 1), nil
		case Renew:
			return 350161 + (period - 1), nil
		case Transfer:
			return 350160, nil
		}
	case "info":
		switch operation{
		case Register:
			return 350051 + (period - 1), nil
		case Renew:
			return 350062 + (period - 1), nil
		case Transfer:
			return 350061, nil
		}
	case "biz":
		switch operation{
		case Register:
			return 350076 + (period - 1), nil
		case Renew:
			return 350087 + (period - 1), nil
		case Transfer:
			return 350086, nil
		}
	case "us":
		switch operation{
		case Register:
			return 350126 + (period - 1), nil
		case Renew:
			return 350137 + (period - 1), nil
		case Transfer:
			return 350136, nil
		}
	}
	return 0, errors.New("Unknow tld or operation")

}

