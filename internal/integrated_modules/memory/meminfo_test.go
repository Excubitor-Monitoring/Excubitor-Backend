package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const input string = `MemTotal:       30500812 kB
MemFree:        21509628 kB
MemAvailable:   25526628 kB
Buffers:           22532 kB
Cached:          4321956 kB
SwapCached:            0 kB
Active:          6801788 kB
Inactive:        1505024 kB
Active(anon):    3980444 kB
Inactive(anon):    21068 kB
Active(file):    2821344 kB
Inactive(file):  1483956 kB
Unevictable:        5748 kB
Mlocked:            5748 kB
SwapTotal:       8388604 kB
SwapFree:        8388604 kB
Zswap:                 0 kB
Zswapped:              0 kB
Dirty:              3776 kB
Writeback:             0 kB
AnonPages:       3968232 kB
Mapped:          1114548 kB
Shmem:             33848 kB
KReclaimable:     121172 kB
Slab:             292964 kB
SReclaimable:     121172 kB
SUnreclaim:       171792 kB
KernelStack:       22384 kB
PageTables:        43892 kB
SecPageTables:         0 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    23639008 kB
Committed_AS:   10730044 kB
VmallocTotal:   34359738367 kB
VmallocUsed:      100636 kB
VmallocChunk:          0 kB
Percpu:            15552 kB
HardwareCorrupted:     0 kB
AnonHugePages:         0 kB
ShmemHugePages:        0 kB
ShmemPmdMapped:        0 kB
FileHugePages:         0 kB
FilePmdMapped:         0 kB
CmaTotal:              0 kB
CmaFree:               0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
Hugetlb:               0 kB
DirectMap4k:      380988 kB
DirectMap2M:    11915264 kB
DirectMap1G:    18874368 kB
`

func TestParseMemInfo(t *testing.T) {
	entries, err := parseMemInfo(input)
	if err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, 30500812, entries["MemTotal"])
	assert.EqualValues(t, 21509628, entries["MemFree"])
	assert.EqualValues(t, 25526628, entries["MemAvailable"])
	assert.EqualValues(t, 22532, entries["Buffers"])
	assert.EqualValues(t, 4321956, entries["Cached"])
	assert.EqualValues(t, 0, entries["SwapCached"])
	assert.EqualValues(t, 6801788, entries["Active"])
	assert.EqualValues(t, 1505024, entries["Inactive"])
	assert.EqualValues(t, 3980444, entries["Active(anon)"])
	assert.EqualValues(t, 21068, entries["Inactive(anon)"])
	assert.EqualValues(t, 2821344, entries["Active(file)"])
	assert.EqualValues(t, 5748, entries["Unevictable"])
	assert.EqualValues(t, 5748, entries["Mlocked"])
	assert.EqualValues(t, 8388604, entries["SwapTotal"])
	assert.EqualValues(t, 8388604, entries["SwapFree"])
	assert.EqualValues(t, 0, entries["Zswap"])
	assert.EqualValues(t, 0, entries["Zswapped"])
	assert.EqualValues(t, 3776, entries["Dirty"])
	assert.EqualValues(t, 0, entries["Writeback"])
	assert.EqualValues(t, 3968232, entries["AnonPages"])
	assert.EqualValues(t, 1114548, entries["Mapped"])
	assert.EqualValues(t, 33848, entries["Shmem"])
	assert.EqualValues(t, 121172, entries["KReclaimable"])
	assert.EqualValues(t, 292964, entries["Slab"])
	assert.EqualValues(t, 121172, entries["SReclaimable"])
	assert.EqualValues(t, 171792, entries["SUnreclaim"])
	assert.EqualValues(t, 22384, entries["KernelStack"])
	assert.EqualValues(t, 43892, entries["PageTables"])
	assert.EqualValues(t, 0, entries["SecPageTables"])
	assert.EqualValues(t, 0, entries["NFS_Unstable"])
	assert.EqualValues(t, 0, entries["Bounce"])
	assert.EqualValues(t, 0, entries["WritebackTmp"])
	assert.EqualValues(t, 23639008, entries["CommitLimit"])
	assert.EqualValues(t, 10730044, entries["Committed_AS"])
	assert.EqualValues(t, 34359738367, entries["VmallocTotal"])
	assert.EqualValues(t, 100636, entries["VmallocUsed"])
	assert.EqualValues(t, 0, entries["VmallocChunk"])
	assert.EqualValues(t, 15552, entries["Percpu"])
	assert.EqualValues(t, 0, entries["HardwareCorrupted"])
	assert.EqualValues(t, 0, entries["AnonHugePages"])
	assert.EqualValues(t, 0, entries["ShmemHugePages"])
	assert.EqualValues(t, 0, entries["ShmemPmdMapped"])
	assert.EqualValues(t, 0, entries["FileHugePages"])
	assert.EqualValues(t, 0, entries["FilePmdMapped"])
	assert.EqualValues(t, 0, entries["CmaTotal"])
	assert.EqualValues(t, 0, entries["CmaFree"])
	assert.EqualValues(t, 0, entries["HugePages_Total"])
	assert.EqualValues(t, 0, entries["HugePages_Free"])
	assert.EqualValues(t, 0, entries["HugePages_Rsvd"])
	assert.EqualValues(t, 0, entries["HugePages_Surp"])
	assert.EqualValues(t, 2048, entries["Hugepagesize"])
	assert.EqualValues(t, 0, entries["Hugetlb"])
	assert.EqualValues(t, 380988, entries["DirectMap4k"])
	assert.EqualValues(t, 11915264, entries["DirectMap2M"])
	assert.EqualValues(t, 18874368, entries["DirectMap1G"])
}

func TestGetMemInfoFromMap(t *testing.T) {
	entries, err := parseMemInfo(input)
	if err != nil {
		t.Error(err)
		return
	}

	memInfo := getMemInfoFromMap(entries)

	assert.EqualValues(t, 30500812, memInfo.MemTotal)
	assert.EqualValues(t, 21509628, memInfo.MemFree)
	assert.EqualValues(t, 25526628, memInfo.MemAvailable)
}

func TestGetSwapInfoFromMap(t *testing.T) {
	entries, err := parseMemInfo(input)
	if err != nil {
		t.Error(err)
		return
	}

	swapInfo := getSwapInfoFromMap(entries)

	assert.EqualValues(t, 8388604, swapInfo.SwapTotal)
	assert.EqualValues(t, 8388604, swapInfo.SwapFree)
}
