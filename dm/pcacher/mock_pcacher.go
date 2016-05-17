/*
   mock_pcacher.go 完成了一个虚拟的cache, 其逻辑和pcacher一模一样, 只不过它将所有的数据
   完全存于内存中.

   mock_pcacher主要用于测试
*/
package pcacher

import (
	"sync"
	"sync/atomic"
)

/*
type Pcacher interface {
    NewPage(initData []byte) (Page, error) // 新创建一页
    GetPage(pgno Pgno) (Page, error)       // 根据叶号取得一页
    Close()
}
*/

type mockPcacher struct {
	cache map[Pgno]*mockPage
	lock  sync.Mutex

	noPages uint32
}

func NewMockPcacher() *mockPcacher {
	return &mockPcacher{
		cache: make(map[Pgno]*mockPage),
	}
}

func (mpc *mockPcacher) NewPage(initData []byte) (Page, error) {
	pgno := Pgno(atomic.AddUint32(&mpc.noPages, 1))
	pg := newMockPage(pgno, initData)
	mpc.lock.Lock()
	defer mpc.lock.Unlock()

	mpc.cache[pgno] = pg
	return pg, nil
}

func (mpc *mockPcacher) GetPage(pgno Pgno) (Page, error) {
	mpc.lock.Lock()
	defer mpc.lock.Unlock()
	return mpc.cache[pgno], nil
}

func (mpc *mockPcacher) Close() {
	// do nothing
}

type mockPage struct {
	pgno Pgno
	data []byte
	lock sync.Mutex
}

func newMockPage(pgno Pgno, data []byte) *mockPage {
	return &mockPage{
		pgno: pgno,
		data: data,
	}
}

func (mp *mockPage) Unlock() {
	mp.lock.Unlock()
}

func (mp *mockPage) Lock() {
	mp.lock.Lock()
}

func (mp *mockPage) Release() {
	// do nothing.
}

func (mp *mockPage) Dirty() {
	// do nothing
}

func (mp *mockPage) Data() []byte {
	return mp.data
}

func (mp *mockPage) Pgno() Pgno {
	return mp.pgno
}