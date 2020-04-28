package ip_guard

import (
	"github.com/Aneg/otus-anti-brute-force/internal/constants"
	"github.com/Aneg/otus-anti-brute-force/internal/models"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories/mock"
	"testing"
)

func TestMemoryIpGuard_AddMask(t *testing.T) {

	rep := &mock.MasksRepository{
		Rows: []models.Mask{
			{Id: 1, Mask: "123.23.44.55/8", ListId: 1},
			{Id: 2, Mask: "122.27.44.55/8", ListId: 1},
		},
	}

	type fields struct {
		listId constants.ListId
		rep    repositories.Masks
	}
	type args struct {
		mask string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "true",
			fields:  fields{listId: 1, rep: rep},
			args:    args{mask: "153.123.44.54/24"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "false",
			fields:  fields{listId: 1, rep: rep},
			args:    args{mask: "123.23.44.55/8"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemoryIpGuard(tt.fields.listId, tt.fields.rep)
			if err := m.Reload(); err != nil {
				t.Error(err)
			}
			got, err := m.AddMask(tt.args.mask)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddMask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddMask() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryIpGuard_DropMask(t *testing.T) {
	rep := &mock.MasksRepository{
		Rows: []models.Mask{
			{Id: 1, Mask: "123.23.44.55/8", ListId: 1},
			{Id: 2, Mask: "122.27.44.55/8", ListId: 1},
		},
	}

	type fields struct {
		listId constants.ListId
		rep    repositories.Masks
	}
	type args struct {
		mask string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "false",
			fields:  fields{listId: 1, rep: rep},
			args:    args{mask: "153.123.44.54/24"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "true",
			fields:  fields{listId: 1, rep: rep},
			args:    args{mask: "123.23.44.55/8"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemoryIpGuard(tt.fields.listId, tt.fields.rep)
			if err := m.Reload(); err != nil {
				t.Error(err)
			}
			got, err := m.DropMask(tt.args.mask)
			if (err != nil) != tt.wantErr {
				t.Errorf("DropMask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DropMask() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryIpGuard_Contains(t *testing.T) {
	rep := &mock.MasksRepository{
		Rows: []models.Mask{
			{Id: 1, Mask: "123.23.44.55/8", ListId: 1},
			{Id: 2, Mask: "122.27.44.55/8", ListId: 1},
		},
	}

	type fields struct {
		listId constants.ListId
		rep    repositories.Masks
	}
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "false",
			fields:  fields{listId: 1, rep: rep},
			args:    args{ip: "123.23.41.55"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "false",
			fields:  fields{listId: 1, rep: rep},
			args:    args{ip: "111.123.44.54/24"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemoryIpGuard(tt.fields.listId, tt.fields.rep)
			if err := m.Reload(); err != nil {
				t.Error(err)
			}
			got, err := m.Contains(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Contains() got = %v, want %v", got, tt.want)
			}
		})
	}
}
