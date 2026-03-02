package lua

func GetSlidingWindowScript() string {
	script := `
	local key = KEYS[1]

	local capacity = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local requested = tonumber(ARGV[4])

	-- Fetch the data from redis
	local result = redis.call("HMGET", key, "currentCnt", "previousCnt", "windowStart")
	local currentCnt = tonumber(result[1]) or 0
	local previousCnt = tonumber(result[2]) or 0
	local windowStart = tonumber(result[3]) or now

	-- Determine the ellapsed time
	local ellapsed = now - windowStart

	-- Shift windows if the current window has passed
	if ellapsed >= window then
		local shift = math.ceil(ellapsed / window)

		if shift >= 2 then
			previousCnt = 0
		else
			previousCnt = currentCnt
		end

		currentCnt = 0
		-- shift the windoStart by however many windows passed
		windowStart = windowStart + (shift * window)
		ellapsed = now - windowStart
	end

	local weight = (window - ellapsed) / window
	local effectiveCnt = currentCnt + ( previousCnt * weight )

	local allowed = 0
	if effectiveCnt + requested <= capacity then
		allowed = 1
		currentCnt = effectiveCnt + requested
	end

	-- Store the data
	redis.call("HMSET", key, "currentCnt", currentCnt, "previousCnt", previousCnt, "windowStart", windowStart)

	-- SET TTL
	local ttl = 2 * window

	redis.call("PEXPIRE", key, ttl)

	return allowed
	`

	return script
}
