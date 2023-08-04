package models

type CreateRequest struct {
	Url        string `json:"url" example:"www.netlify.com"`
	ExpireTime int    `json:"expireTime" example:"10"`
	CustomTag  string `json:"customTag" example:"net-short"`
}

func (req CreateRequest) EnforceHttp() {
	if req.Url[:4] != "http" {
		req.Url = "https://" + req.Url
	}
}

func (req CreateRequest) DefaultExpireTime() {
	if req.ExpireTime == 0 {
		req.ExpireTime = 24
	}
}

type CreateResponse struct {
	Url             string `json:"url" example:"www.netlify.com"`
	CustomTag       string `json:"customTag" example:"net-short"`
	ExpireTime      int    `json:"expireTime" example:"10"`
	XRateRemaining  int    `json:"rateLimit" example:"29"`
	XRateLimitReset int    `json:"rateLimitReset" example:"29"`
}

type BadRequestError struct {
	Error       string `json:"error" example:"You exhausted your limit"`
	TimeToReset int    `json:"timeToReset" example:"10"`
	StatusCode  int    `json:"statusCode" example:"400"`
}

type RequestError struct {
	Error      string `json:"error" example:"Internal Server Error"`
	StatusCode int    `json:"statusCode" example:"500"`
}
