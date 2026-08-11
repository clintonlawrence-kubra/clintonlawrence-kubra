package rbac_test

import data.rbac

test_admin_bypasses_quota if {
	rbac.allow with input as {"user": {"role": "admin", "plan": "free"}, "requests_today": 999}
}

test_free_under_limit_allowed if {
	rbac.allow with input as {"user": {"role": "user", "plan": "free"}, "requests_today": 2}
}

test_free_at_limit_denied if {
	not rbac.allow with input as {"user": {"role": "user", "plan": "free"}, "requests_today": 3}
}

test_free_over_limit_denied if {
	not rbac.allow with input as {"user": {"role": "user", "plan": "free"}, "requests_today": 5}
}

test_pro_under_limit_allowed if {
	rbac.allow with input as {"user": {"role": "user", "plan": "pro"}, "requests_today": 99}
}

test_pro_at_limit_denied if {
	not rbac.allow with input as {"user": {"role": "user", "plan": "pro"}, "requests_today": 100}
}

test_unknown_plan_denied if {
	not rbac.allow with input as {"user": {"role": "user", "plan": "enterprise"}, "requests_today": 0}
}
