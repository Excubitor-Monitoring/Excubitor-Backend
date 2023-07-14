package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStat(t *testing.T) {
	input := `cpu  30739 199 8810 797183 1320 1548 735 0 0 0
cpu0 2221 19 623 49260 82 146 175 0 0 0
cpu1 1260 7 447 50564 65 110 83 0 0 0
cpu2 2139 23 598 49517 115 94 42 0 0 0
cpu3 1773 47 421 50129 97 56 22 0 0 0
cpu4 2295 7 595 49406 92 92 32 0 0 0
cpu5 1545 14 853 49788 51 245 19 0 0 0
cpu6 1971 3 654 49718 58 88 29 0 0 0
cpu7 1890 7 477 50000 61 66 21 0 0 0
cpu8 2567 6 569 49028 127 112 101 0 0 0
cpu9 1374 1 427 50515 89 73 62 0 0 0
cpu10 1956 6 623 49753 74 93 31 0 0 0
cpu11 1248 5 392 50724 95 60 18 0 0 0
cpu12 1880 14 572 49859 77 94 30 0 0 0
cpu13 2305 5 484 49602 65 66 19 0 0 0
cpu14 2149 9 607 49568 78 82 27 0 0 0
cpu15 2157 18 458 49746 87 63 16 0 0 0
intr 2696648 118 1821 0 0 0 0 0 0 1 13719 0 0 40316 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 297 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 21994 0 0 0 0 0 0 0 0 0 0 40 4935 2570 4495 3175 3417 2949 4635 2617 4468 7948 5334 3782 3517 3419 6151 3014 81337 0 0 1 343 1296 0 60498 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 2883660
btime 1684405163
processes 9548
procs_running 1
procs_blocked 1
softirq 1157133 13604 76082 4 25454 46 0 178231 408577 151 454984`

	stats, err := parseStat(input)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 17, len(stats))
	assert.Equal(t, "cpu", stats[0].name)
	assert.Equal(t, "cpu0", stats[1].name)
	assert.Equal(t, "cpu1", stats[2].name)
	assert.Equal(t, "cpu2", stats[3].name)
	assert.Equal(t, "cpu3", stats[4].name)
	assert.Equal(t, "cpu4", stats[5].name)
	assert.Equal(t, "cpu5", stats[6].name)
	assert.Equal(t, "cpu6", stats[7].name)
	assert.Equal(t, "cpu7", stats[8].name)
	assert.Equal(t, "cpu8", stats[9].name)
	assert.Equal(t, "cpu9", stats[10].name)
	assert.Equal(t, "cpu10", stats[11].name)
	assert.Equal(t, "cpu11", stats[12].name)
	assert.Equal(t, "cpu12", stats[13].name)
	assert.Equal(t, "cpu13", stats[14].name)
	assert.Equal(t, "cpu14", stats[15].name)
	assert.Equal(t, "cpu15", stats[16].name)
}

func TestParseCoreStat(t *testing.T) {
	line := "cpu0 2221 19 623 49260 82 146 175 0 0 0\n"

	stat, err := parseCoreStat(line)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "cpu0", stat.name)
	assert.EqualValues(t, 2221, stat.User)
	assert.EqualValues(t, 19, stat.Nice)
	assert.EqualValues(t, 623, stat.System)
	assert.EqualValues(t, 49260, stat.Idle)
	assert.EqualValues(t, 82, stat.Iowait)
	assert.EqualValues(t, 146, stat.Irq)
	assert.EqualValues(t, 175, stat.Softirq)
	assert.EqualValues(t, 0, stat.Steal)
	assert.EqualValues(t, 0, stat.Guest)
}
