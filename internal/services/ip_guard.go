package services

type IpGuard interface {
	AddMask(mask string) (bool, error)  // добавляем если её нету в списке
	DropMask(mask string) (bool, error) // удаляем если есть в списке
	Contains(ip string) (bool, error)   // проверяем на соответствие имеющимся маскам ip
	Reload() error
}
