package codemessage

const activationHTMLTpl = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
</head>
<body style="margin:0;padding:0;background-color:#f0f2f5;">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="background-color:#f0f2f5;">
<tr><td align="center" style="padding:48px 20px;">
<table role="presentation" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:20px;overflow:hidden;">

	<tr>
		<td align="center" style="background:linear-gradient(135deg,#7c3aed 0%,#5b21b6 100%);padding:44px 48px 40px;">
			<p style="margin:0 0 6px;color:#ffffff;font-size:26px;font-weight:700;font-family:Arial,Helvetica,sans-serif;">Messenger</p>
		<p style="margin:0;color:rgba(255,255,255,0.82);font-size:13px;font-weight:600;letter-spacing:1.8px;text-transform:uppercase;font-family:Arial,Helvetica,sans-serif;">{SUBTITLE}</p>
		</td>
	</tr>

	<tr>
		<td style="padding:40px 48px 36px;">
			<p style="margin:0 0 32px;color:#374151;font-size:15px;line-height:1.75;font-family:Arial,Helvetica,sans-serif;">{BODY}</p>

			<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="margin-bottom:32px;">
				<tr>
					<td align="center" style="background-color:#f5f3ff;border:2px solid #ede9fe;border-radius:16px;padding:28px 20px;">
						<p style="margin:0 0 10px;color:#7c3aed;font-size:11px;font-weight:700;letter-spacing:2.5px;text-transform:uppercase;font-family:Arial,Helvetica,sans-serif;">{CODE_LABEL}</p>
						<p style="margin:0;color:#4c1d95;font-size:46px;font-weight:800;letter-spacing:16px;padding-left:16px;font-family:'Courier New',Courier,monospace;">{CODE}</p>
					</td>
				</tr>
			</table>

			<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="margin-bottom:32px;">
				<tr>
					<td style="border-top:1px solid #e5e7eb;"></td>
					<td style="padding:0 14px;white-space:nowrap;color:#9ca3af;font-size:13px;font-family:Arial,Helvetica,sans-serif;">{DIVIDER}</td>
					<td style="border-top:1px solid #e5e7eb;"></td>
				</tr>
			</table>

			<table role="presentation" width="100%" cellpadding="0" cellspacing="0">
				<tr>
					<td align="center">
						<a href="{TOKEN_URL}" style="display:inline-block;background:linear-gradient(135deg,#7c3aed 0%,#5b21b6 100%);color:#ffffff;text-decoration:none;font-size:15px;font-weight:700;padding:16px 56px;border-radius:50px;font-family:Arial,Helvetica,sans-serif;">{BUTTON}</a>
					</td>
				</tr>
			</table>
		</td>
	</tr>

	<tr>
		<td style="background-color:#f9fafb;border-top:1px solid #f3f4f6;padding:24px 48px;">
			<p style="margin:0;color:#9ca3af;font-size:12px;line-height:1.8;text-align:center;font-family:Arial,Helvetica,sans-serif;">{FOOTER}</p>
		</td>
	</tr>

</table>
</td></tr>
</table>
</body>
</html>`

const passwordResetHTMLTpl = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
</head>
<body style="margin:0;padding:0;background-color:#f0f2f5;">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="background-color:#f0f2f5;">
<tr><td align="center" style="padding:48px 20px;">
<table role="presentation" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;background-color:#ffffff;border-radius:20px;overflow:hidden;">

	<tr>
		<td align="center" style="background:linear-gradient(135deg,#dc2626 0%,#991b1b 100%);padding:44px 48px 40px;">
			<p style="margin:0 0 6px;color:#ffffff;font-size:26px;font-weight:700;font-family:Arial,Helvetica,sans-serif;">Messenger</p>
			<p style="margin:0;color:rgba(255,255,255,0.82);font-size:13px;font-weight:600;letter-spacing:1.8px;text-transform:uppercase;font-family:Arial,Helvetica,sans-serif;">{SUBTITLE}</p>
		</td>
	</tr>

	<tr>
		<td style="background-color:#fef2f2;border-bottom:1px solid #fee2e2;padding:14px 48px;">
			<p style="margin:0;color:#991b1b;font-size:13px;text-align:center;font-family:Arial,Helvetica,sans-serif;">&#9888;&#65039; {WARNING}</p>
		</td>
	</tr>

	<tr>
		<td style="padding:40px 48px 36px;">
			<p style="margin:0 0 32px;color:#374151;font-size:15px;line-height:1.75;font-family:Arial,Helvetica,sans-serif;">{BODY}</p>

			<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="margin-bottom:32px;">
				<tr>
					<td align="center" style="background-color:#fff7f7;border:2px solid #fecaca;border-radius:16px;padding:28px 20px;">
						<p style="margin:0 0 10px;color:#dc2626;font-size:11px;font-weight:700;letter-spacing:2.5px;text-transform:uppercase;font-family:Arial,Helvetica,sans-serif;">{CODE_LABEL}</p>
						<p style="margin:0;color:#7f1d1d;font-size:46px;font-weight:800;letter-spacing:16px;padding-left:16px;font-family:'Courier New',Courier,monospace;">{CODE}</p>
					</td>
				</tr>
			</table>

			<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="margin-bottom:32px;">
				<tr>
					<td style="border-top:1px solid #e5e7eb;"></td>
					<td style="padding:0 14px;white-space:nowrap;color:#9ca3af;font-size:13px;font-family:Arial,Helvetica,sans-serif;">{DIVIDER}</td>
					<td style="border-top:1px solid #e5e7eb;"></td>
				</tr>
			</table>

			<table role="presentation" width="100%" cellpadding="0" cellspacing="0">
				<tr>
					<td align="center">
						<a href="{TOKEN_URL}" style="display:inline-block;background:linear-gradient(135deg,#dc2626 0%,#991b1b 100%);color:#ffffff;text-decoration:none;font-size:15px;font-weight:700;padding:16px 56px;border-radius:50px;font-family:Arial,Helvetica,sans-serif;">{BUTTON}</a>
					</td>
				</tr>
			</table>
		</td>
	</tr>

	<tr>
		<td style="background-color:#f9fafb;border-top:1px solid #f3f4f6;padding:24px 48px;">
			<p style="margin:0;color:#9ca3af;font-size:12px;line-height:1.8;text-align:center;font-family:Arial,Helvetica,sans-serif;">{FOOTER}</p>
		</td>
	</tr>

</table>
</td></tr>
</table>
</body>
</html>`