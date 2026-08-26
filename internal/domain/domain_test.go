package domain_test

import (
	"errors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"testing"
	"time"
)

func TestResourceTransitions(t *testing.T) {
	cases := []struct {
		from, to domain.ResourceStatus
		ok       bool
	}{{domain.ResourceDraft, domain.ResourceSubmitted, true}, {domain.ResourceSubmitted, domain.ResourceReviewing, true}, {domain.ResourceReviewing, domain.ResourceRegistered, true}, {domain.ResourceRegistered, domain.ResourcePublished, true}, {domain.ResourceDraft, domain.ResourcePublished, false}}
	for _, c := range cases {
		if c.from.Can(c.to) != c.ok {
			t.Fatal(c)
		}
	}
}
func TestValidation(t *testing.T) {
	now := time.Now()
	r := domain.Resource{ID: "r", TenantID: "t", OwnerID: "u", Code: "gansu-01", Name: "人口", CreatedAt: now}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
	r.Name = " "
	if !errors.Is(r.Validate(), domain.ErrEmptyName) {
		t.Fatal()
	}
}
func TestResourceValidationVariant0(t *testing.T) {
	r := domain.Resource{ID: "r0", TenantID: "t", OwnerID: "u", Code: "code-0", Name: "名称0"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant1(t *testing.T) {
	r := domain.Resource{ID: "r1", TenantID: "t", OwnerID: "u", Code: "code-1", Name: "名称1"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant2(t *testing.T) {
	r := domain.Resource{ID: "r2", TenantID: "t", OwnerID: "u", Code: "code-2", Name: "名称2"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant3(t *testing.T) {
	r := domain.Resource{ID: "r3", TenantID: "t", OwnerID: "u", Code: "code-3", Name: "名称3"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant4(t *testing.T) {
	r := domain.Resource{ID: "r4", TenantID: "t", OwnerID: "u", Code: "code-4", Name: "名称4"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant5(t *testing.T) {
	r := domain.Resource{ID: "r5", TenantID: "t", OwnerID: "u", Code: "code-5", Name: "名称5"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant6(t *testing.T) {
	r := domain.Resource{ID: "r6", TenantID: "t", OwnerID: "u", Code: "code-6", Name: "名称6"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant7(t *testing.T) {
	r := domain.Resource{ID: "r7", TenantID: "t", OwnerID: "u", Code: "code-7", Name: "名称7"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant8(t *testing.T) {
	r := domain.Resource{ID: "r8", TenantID: "t", OwnerID: "u", Code: "code-8", Name: "名称8"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant9(t *testing.T) {
	r := domain.Resource{ID: "r9", TenantID: "t", OwnerID: "u", Code: "code-9", Name: "名称9"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant10(t *testing.T) {
	r := domain.Resource{ID: "r10", TenantID: "t", OwnerID: "u", Code: "code-10", Name: "名称10"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant11(t *testing.T) {
	r := domain.Resource{ID: "r11", TenantID: "t", OwnerID: "u", Code: "code-11", Name: "名称11"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant12(t *testing.T) {
	r := domain.Resource{ID: "r12", TenantID: "t", OwnerID: "u", Code: "code-12", Name: "名称12"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant13(t *testing.T) {
	r := domain.Resource{ID: "r13", TenantID: "t", OwnerID: "u", Code: "code-13", Name: "名称13"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant14(t *testing.T) {
	r := domain.Resource{ID: "r14", TenantID: "t", OwnerID: "u", Code: "code-14", Name: "名称14"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant15(t *testing.T) {
	r := domain.Resource{ID: "r15", TenantID: "t", OwnerID: "u", Code: "code-15", Name: "名称15"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant16(t *testing.T) {
	r := domain.Resource{ID: "r16", TenantID: "t", OwnerID: "u", Code: "code-16", Name: "名称16"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant17(t *testing.T) {
	r := domain.Resource{ID: "r17", TenantID: "t", OwnerID: "u", Code: "code-17", Name: "名称17"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant18(t *testing.T) {
	r := domain.Resource{ID: "r18", TenantID: "t", OwnerID: "u", Code: "code-18", Name: "名称18"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant19(t *testing.T) {
	r := domain.Resource{ID: "r19", TenantID: "t", OwnerID: "u", Code: "code-19", Name: "名称19"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant20(t *testing.T) {
	r := domain.Resource{ID: "r20", TenantID: "t", OwnerID: "u", Code: "code-20", Name: "名称20"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant21(t *testing.T) {
	r := domain.Resource{ID: "r21", TenantID: "t", OwnerID: "u", Code: "code-21", Name: "名称21"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant22(t *testing.T) {
	r := domain.Resource{ID: "r22", TenantID: "t", OwnerID: "u", Code: "code-22", Name: "名称22"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant23(t *testing.T) {
	r := domain.Resource{ID: "r23", TenantID: "t", OwnerID: "u", Code: "code-23", Name: "名称23"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant24(t *testing.T) {
	r := domain.Resource{ID: "r24", TenantID: "t", OwnerID: "u", Code: "code-24", Name: "名称24"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant25(t *testing.T) {
	r := domain.Resource{ID: "r25", TenantID: "t", OwnerID: "u", Code: "code-25", Name: "名称25"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant26(t *testing.T) {
	r := domain.Resource{ID: "r26", TenantID: "t", OwnerID: "u", Code: "code-26", Name: "名称26"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant27(t *testing.T) {
	r := domain.Resource{ID: "r27", TenantID: "t", OwnerID: "u", Code: "code-27", Name: "名称27"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant28(t *testing.T) {
	r := domain.Resource{ID: "r28", TenantID: "t", OwnerID: "u", Code: "code-28", Name: "名称28"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant29(t *testing.T) {
	r := domain.Resource{ID: "r29", TenantID: "t", OwnerID: "u", Code: "code-29", Name: "名称29"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant30(t *testing.T) {
	r := domain.Resource{ID: "r30", TenantID: "t", OwnerID: "u", Code: "code-30", Name: "名称30"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant31(t *testing.T) {
	r := domain.Resource{ID: "r31", TenantID: "t", OwnerID: "u", Code: "code-31", Name: "名称31"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant32(t *testing.T) {
	r := domain.Resource{ID: "r32", TenantID: "t", OwnerID: "u", Code: "code-32", Name: "名称32"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant33(t *testing.T) {
	r := domain.Resource{ID: "r33", TenantID: "t", OwnerID: "u", Code: "code-33", Name: "名称33"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant34(t *testing.T) {
	r := domain.Resource{ID: "r34", TenantID: "t", OwnerID: "u", Code: "code-34", Name: "名称34"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant35(t *testing.T) {
	r := domain.Resource{ID: "r35", TenantID: "t", OwnerID: "u", Code: "code-35", Name: "名称35"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant36(t *testing.T) {
	r := domain.Resource{ID: "r36", TenantID: "t", OwnerID: "u", Code: "code-36", Name: "名称36"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant37(t *testing.T) {
	r := domain.Resource{ID: "r37", TenantID: "t", OwnerID: "u", Code: "code-37", Name: "名称37"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant38(t *testing.T) {
	r := domain.Resource{ID: "r38", TenantID: "t", OwnerID: "u", Code: "code-38", Name: "名称38"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant39(t *testing.T) {
	r := domain.Resource{ID: "r39", TenantID: "t", OwnerID: "u", Code: "code-39", Name: "名称39"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant40(t *testing.T) {
	r := domain.Resource{ID: "r40", TenantID: "t", OwnerID: "u", Code: "code-40", Name: "名称40"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant41(t *testing.T) {
	r := domain.Resource{ID: "r41", TenantID: "t", OwnerID: "u", Code: "code-41", Name: "名称41"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant42(t *testing.T) {
	r := domain.Resource{ID: "r42", TenantID: "t", OwnerID: "u", Code: "code-42", Name: "名称42"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant43(t *testing.T) {
	r := domain.Resource{ID: "r43", TenantID: "t", OwnerID: "u", Code: "code-43", Name: "名称43"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant44(t *testing.T) {
	r := domain.Resource{ID: "r44", TenantID: "t", OwnerID: "u", Code: "code-44", Name: "名称44"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant45(t *testing.T) {
	r := domain.Resource{ID: "r45", TenantID: "t", OwnerID: "u", Code: "code-45", Name: "名称45"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant46(t *testing.T) {
	r := domain.Resource{ID: "r46", TenantID: "t", OwnerID: "u", Code: "code-46", Name: "名称46"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant47(t *testing.T) {
	r := domain.Resource{ID: "r47", TenantID: "t", OwnerID: "u", Code: "code-47", Name: "名称47"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant48(t *testing.T) {
	r := domain.Resource{ID: "r48", TenantID: "t", OwnerID: "u", Code: "code-48", Name: "名称48"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant49(t *testing.T) {
	r := domain.Resource{ID: "r49", TenantID: "t", OwnerID: "u", Code: "code-49", Name: "名称49"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant50(t *testing.T) {
	r := domain.Resource{ID: "r50", TenantID: "t", OwnerID: "u", Code: "code-50", Name: "名称50"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant51(t *testing.T) {
	r := domain.Resource{ID: "r51", TenantID: "t", OwnerID: "u", Code: "code-51", Name: "名称51"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant52(t *testing.T) {
	r := domain.Resource{ID: "r52", TenantID: "t", OwnerID: "u", Code: "code-52", Name: "名称52"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant53(t *testing.T) {
	r := domain.Resource{ID: "r53", TenantID: "t", OwnerID: "u", Code: "code-53", Name: "名称53"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant54(t *testing.T) {
	r := domain.Resource{ID: "r54", TenantID: "t", OwnerID: "u", Code: "code-54", Name: "名称54"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant55(t *testing.T) {
	r := domain.Resource{ID: "r55", TenantID: "t", OwnerID: "u", Code: "code-55", Name: "名称55"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant56(t *testing.T) {
	r := domain.Resource{ID: "r56", TenantID: "t", OwnerID: "u", Code: "code-56", Name: "名称56"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant57(t *testing.T) {
	r := domain.Resource{ID: "r57", TenantID: "t", OwnerID: "u", Code: "code-57", Name: "名称57"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant58(t *testing.T) {
	r := domain.Resource{ID: "r58", TenantID: "t", OwnerID: "u", Code: "code-58", Name: "名称58"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant59(t *testing.T) {
	r := domain.Resource{ID: "r59", TenantID: "t", OwnerID: "u", Code: "code-59", Name: "名称59"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant60(t *testing.T) {
	r := domain.Resource{ID: "r60", TenantID: "t", OwnerID: "u", Code: "code-60", Name: "名称60"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant61(t *testing.T) {
	r := domain.Resource{ID: "r61", TenantID: "t", OwnerID: "u", Code: "code-61", Name: "名称61"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant62(t *testing.T) {
	r := domain.Resource{ID: "r62", TenantID: "t", OwnerID: "u", Code: "code-62", Name: "名称62"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant63(t *testing.T) {
	r := domain.Resource{ID: "r63", TenantID: "t", OwnerID: "u", Code: "code-63", Name: "名称63"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant64(t *testing.T) {
	r := domain.Resource{ID: "r64", TenantID: "t", OwnerID: "u", Code: "code-64", Name: "名称64"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant65(t *testing.T) {
	r := domain.Resource{ID: "r65", TenantID: "t", OwnerID: "u", Code: "code-65", Name: "名称65"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant66(t *testing.T) {
	r := domain.Resource{ID: "r66", TenantID: "t", OwnerID: "u", Code: "code-66", Name: "名称66"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant67(t *testing.T) {
	r := domain.Resource{ID: "r67", TenantID: "t", OwnerID: "u", Code: "code-67", Name: "名称67"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant68(t *testing.T) {
	r := domain.Resource{ID: "r68", TenantID: "t", OwnerID: "u", Code: "code-68", Name: "名称68"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant69(t *testing.T) {
	r := domain.Resource{ID: "r69", TenantID: "t", OwnerID: "u", Code: "code-69", Name: "名称69"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant70(t *testing.T) {
	r := domain.Resource{ID: "r70", TenantID: "t", OwnerID: "u", Code: "code-70", Name: "名称70"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant71(t *testing.T) {
	r := domain.Resource{ID: "r71", TenantID: "t", OwnerID: "u", Code: "code-71", Name: "名称71"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant72(t *testing.T) {
	r := domain.Resource{ID: "r72", TenantID: "t", OwnerID: "u", Code: "code-72", Name: "名称72"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant73(t *testing.T) {
	r := domain.Resource{ID: "r73", TenantID: "t", OwnerID: "u", Code: "code-73", Name: "名称73"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant74(t *testing.T) {
	r := domain.Resource{ID: "r74", TenantID: "t", OwnerID: "u", Code: "code-74", Name: "名称74"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant75(t *testing.T) {
	r := domain.Resource{ID: "r75", TenantID: "t", OwnerID: "u", Code: "code-75", Name: "名称75"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant76(t *testing.T) {
	r := domain.Resource{ID: "r76", TenantID: "t", OwnerID: "u", Code: "code-76", Name: "名称76"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant77(t *testing.T) {
	r := domain.Resource{ID: "r77", TenantID: "t", OwnerID: "u", Code: "code-77", Name: "名称77"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant78(t *testing.T) {
	r := domain.Resource{ID: "r78", TenantID: "t", OwnerID: "u", Code: "code-78", Name: "名称78"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestResourceValidationVariant79(t *testing.T) {
	r := domain.Resource{ID: "r79", TenantID: "t", OwnerID: "u", Code: "code-79", Name: "名称79"}
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
