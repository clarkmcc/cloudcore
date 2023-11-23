package agent

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestIsExpiringSoon(t *testing.T) {
	now := time.Unix(0, 0)
	tests := map[string]struct {
		now  time.Time
		exp  time.Time
		dur  time.Duration
		want bool
	}{
		"Less than 30% life remaining": {
			now:  now,
			exp:  now.Add(3 * time.Hour),
			dur:  10 * time.Hour,
			want: true,
		},
		"Exactly 30% life remaining": {
			now:  now,
			exp:  now.Add(7 * time.Hour),
			dur:  10 * time.Hour,
			want: false,
		},
		"More than 30% life remaining": {
			now:  now,
			exp:  now.Add(6 * time.Hour),
			dur:  10 * time.Hour,
			want: false,
		},
		"Expired token": {
			now:  now.Add(12 * time.Hour),
			exp:  now.Add(10 * time.Hour),
			dur:  10 * time.Hour,
			want: true,
		},
		"Long duration, expiring soon": {
			now:  now.Add(100 * time.Hour),
			exp:  now.Add(150 * time.Hour),
			dur:  150 * time.Hour,
			want: true,
		},
		"Long duration, not expiring soon": {
			now:  now.Add(90 * time.Hour),
			exp:  now.Add(150 * time.Hour),
			dur:  150 * time.Hour,
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, isExpiringSoon(test.now, test.exp, test.dur), test.want)
		})
	}
}
