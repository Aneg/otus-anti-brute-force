package services

type IpGuard interface {
	AddMask(mask string) (bool, error)  // добавляем если её нету в списке
	DropMask(mask string) (bool, error) // удаляем если есть в списке
	Contains(ip string) (bool, error)   // проверяем на соответствие имеющимся маскам ip
	Reload(masks []string) error        // TODO: стоит ли лоадер выносить за логику листа? Тогда можно будет сделать 1 воркер для обновления всех листов.
}
