package handlers

import "context"

// UserIDFromContext extracts the authenticated user id from context if present.
func UserIDFromContext(ctx context.Context) (int64, bool) {
    v := ctx.Value(userIDContextKey)
    if v == nil {
        return 0, false
    }
    id, ok := v.(int64)
    return id, ok
}


