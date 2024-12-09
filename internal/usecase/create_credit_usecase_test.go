package usecase

import (
	"testing"
)

func Test_calculateMonthlyPayment(t *testing.T) {
	type args struct {
		principal          int64
		annualInterestRate int32
		numPayments        int
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "calculateMonthlyPayment",
			args: args{
				principal:          2000000,
				annualInterestRate: 12,
				numPayments:        36,
			},
			want: 66429,
		},
		{
			name: "calculateMonthlyPayment",
			args: args{
				principal:          150000,
				annualInterestRate: 20,
				numPayments:        12,
			},
			want: 13895,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMonthlyPayment(tt.args.principal, tt.args.annualInterestRate, tt.args.numPayments); got != tt.want {
				t.Errorf("calculateMonthlyPayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
