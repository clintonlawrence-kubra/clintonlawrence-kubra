package rbac

default allow := false
quota := {
    "free": 3,
    "pro": 100,
}

actions := {
    "get": ["admin", "user","finance","viewer"],
    "approve": ["admin","finance"],
    "delete": ["admin","finance"],
}

under_quota(plan, requests_today) if requests_today < quota[plan]

allow if input.user.role == "admin"


deny contains "action not allowed" if not input.action in actions[input.user.role] 
deny contains "quota exceeded" if not under_quota(input.user.plan, input.user.requests_today)
deny contains "cant improve its own invoice" if input.user.id == input.invoice.owner_id

allow := count(deny) == 0