package account

const testVersion = 1

type Account struct {
	balance int64
	closed  bool
	lock    chan bool
}

func Open(initialDeposit int64) *Account {
	if initialDeposit < 0 {
		return nil
	}
	lock := make(chan bool, 1)
	lock <- false
	return &Account{initialDeposit, false, lock}
}

func (a *Account) Close() (payout int64, ok bool) {
	<-a.lock
	defer a.unlock()

	if a.closed {
		return 0, false
	}
	a.closed = true
	return a.balance, a.closed
}

func (a *Account) Balance() (balance int64, ok bool) {
	if a.closed {
		return 0, false
	}
	return a.balance, true
}

func (a *Account) Deposit(amount int64) (newBalance int64, ok bool) {
	if a.closed {
		return 0, false
	}

	bal := a.balance + amount
	if bal < 0 {
		return 0, false
	} else {
		a.balance = bal
	}

	return a.balance, true
}

func (a *Account) unlock() {
	a.lock <- false
}
