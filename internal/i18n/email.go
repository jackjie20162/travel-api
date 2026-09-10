package i18n

import "fmt"

// This file holds locale-aware email/notification templates. Per the 多语言多货币
// 改造方案 (P3), the verification-code email is localized now, and order/refund
// notification templates are established as ready-to-use keys.

// EmailVerificationCode renders the localized verification-code email.
// It returns the subject and HTML body for the given locale.
func EmailVerificationCode(locale, code string, expireSeconds int) (subject, body string) {
	switch NormalizeLocale(locale) {
	case LocaleZhCN:
		subject = fmt.Sprintf("验证码：%s", code)
		body = fmt.Sprintf(`<div dir="ltr" style="font-family:sans-serif;padding:20px">
		<h2>您的验证码</h2>
		<p>验证码: <strong style="font-size:24px;color:#e74c3c">%s</strong></p>
		<p>有效期 %d 秒，请勿泄露给他人。</p>
	</div>`, code, expireSeconds)
	case LocaleAr:
		subject = fmt.Sprintf("رمز التحقق: %s", code)
		body = fmt.Sprintf(`<div dir="rtl" style="font-family:sans-serif;padding:20px">
		<h2>رمز التحقق الخاص بك</h2>
		<p>رمز التحقق: <strong style="font-size:24px;color:#e74c3c">%s</strong></p>
		<p>صالح لمدة %d ثانية، يرجى عدم مشاركته مع الآخرين.</p>
	</div>`, code, expireSeconds)
	default: // en-US
		subject = fmt.Sprintf("Verification code: %s", code)
		body = fmt.Sprintf(`<div dir="ltr" style="font-family:sans-serif;padding:20px">
		<h2>Your verification code</h2>
		<p>Code: <strong style="font-size:24px;color:#e74c3c">%s</strong></p>
		<p>Valid for %d seconds. Please do not share it with anyone.</p>
	</div>`, code, expireSeconds)
	}
	return subject, body
}

// EmailOrderConfirmation renders the localized order confirmation email.
// Established per P3 (order notification template key); ready for wiring when the
// order-confirmation email flow is enabled.
func EmailOrderConfirmation(locale, orderNo, amount, currency string) (subject, body string) {
	switch NormalizeLocale(locale) {
	case LocaleZhCN:
		subject = fmt.Sprintf("订单确认：%s", orderNo)
		body = fmt.Sprintf(`<div dir="ltr" style="font-family:sans-serif;padding:20px">
		<h2>订单已确认</h2>
		<p>您的订单 <strong>%s</strong> 已确认。</p>
		<p>金额：%s %s</p>
	</div>`, orderNo, amount, currency)
	case LocaleAr:
		subject = fmt.Sprintf("تأكيد الطلب: %s", orderNo)
		body = fmt.Sprintf(`<div dir="rtl" style="font-family:sans-serif;padding:20px">
		<h2>تم تأكيد الطلب</h2>
		<p>تم تأكيد طلبك <strong>%s</strong>.</p>
		<p>المبلغ: %s %s</p>
	</div>`, orderNo, amount, currency)
	default: // en-US
		subject = fmt.Sprintf("Order confirmation: %s", orderNo)
		body = fmt.Sprintf(`<div dir="ltr" style="font-family:sans-serif;padding:20px">
		<h2>Order confirmed</h2>
		<p>Your order <strong>%s</strong> has been confirmed.</p>
		<p>Amount: %s %s</p>
	</div>`, orderNo, amount, currency)
	}
	return subject, body
}

// EmailRefundNotification renders the localized refund notification email.
// Established per P3 (refund notification template key); ready for wiring when the
// refund-notification email flow is enabled.
func EmailRefundNotification(locale, orderNo, status string) (subject, body string) {
	switch NormalizeLocale(locale) {
	case LocaleZhCN:
		subject = fmt.Sprintf("退款通知：%s", orderNo)
		body = fmt.Sprintf(`<div dir="ltr" style="font-family:sans-serif;padding:20px">
		<h2>退款状态更新</h2>
		<p>您的订单 <strong>%s</strong> 退款状态：%s</p>
	</div>`, orderNo, status)
	case LocaleAr:
		subject = fmt.Sprintf("إشعار الاسترداد: %s", orderNo)
		body = fmt.Sprintf(`<div dir="rtl" style="font-family:sans-serif;padding:20px">
		<h2>تحديث حالة الاسترداد</h2>
		<p>حالة استرداد طلبك <strong>%s</strong>: %s</p>
	</div>`, orderNo, status)
	default: // en-US
		subject = fmt.Sprintf("Refund notification: %s", orderNo)
		body = fmt.Sprintf(`<div dir="ltr" style="font-family:sans-serif;padding:20px">
		<h2>Refund status update</h2>
		<p>Refund status for your order <strong>%s</strong>: %s</p>
	</div>`, orderNo, status)
	}
	return subject, body
}
