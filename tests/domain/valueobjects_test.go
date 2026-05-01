package domain_test

import (
	"testing"

	"github.com/casatorino/backend/internal/domain/valueobjects"
)

func TestValueObjectsValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		run     func() error
		wantErr bool
	}{
		{
			name: "valid customer type",
			run: func() error {
				_, err := valueobjects.NewCustomerType("PERSON")
				return err
			},
		},
		{
			name: "invalid customer type",
			run: func() error {
				_, err := valueobjects.NewCustomerType("UNKNOWN")
				return err
			},
			wantErr: true,
		},
		{
			name: "valid product type",
			run: func() error {
				_, err := valueobjects.NewProductType("CAKE")
				return err
			},
		},
		{
			name: "valid payment method",
			run: func() error {
				_, err := valueobjects.NewPaymentMethod("CARD")
				return err
			},
		},
		{
			name: "invalid payment status",
			run: func() error {
				_, err := valueobjects.NewPaymentStatus("UNKNOWN")
				return err
			},
			wantErr: true,
		},
		{
			name: "valid unit",
			run: func() error {
				_, err := valueobjects.NewUnit("KG")
				return err
			},
		},
		{
			name: "invalid order status",
			run: func() error {
				_, err := valueobjects.NewOrderStatus("ARCHIVED")
				return err
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.run()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
