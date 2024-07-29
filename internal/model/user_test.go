package model

import "testing"

func Test_maskPhoneNumber(t *testing.T) {
	type args struct {
		user User
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "휴대전화 마스킹",
			args: args{
				user: User{
					PhoneNumber: "01055551111",
				},
			},
			want: "0105**51**1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.user.MaskPhoneNumber(); got != tt.want {
				t.Errorf("maskPhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
