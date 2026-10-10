package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type App struct {
	ctx context.Context
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// IPInfo уходит во фронтенд как JSON.
type IPInfo struct {
	IP          string  `json:"ip"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Zip         string  `json:"zip"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	HasGeo      bool    `json:"hasGeo"`   // есть ли координаты для карты
	PingMs      int     `json:"pingMs"`   // время ответа сервиса
	Provider    string  `json:"provider"` // какой сервис ответил
}

var httpClient = &http.Client{Timeout: 7 * time.Second}

func getJSON(url string, v any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) IP-Checker/1.0")
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("bad response: %.40s", strings.TrimSpace(string(body)))
	}
	return nil
}

// geo заполняет координаты и отмечает, что они валидны.
func (i *IPInfo) setGeo(lat, lon float64) {
	if lat == 0 && lon == 0 {
		return
	}
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return
	}
	i.Lat, i.Lon, i.HasGeo = lat, lon, true
}

func num(n json.Number) float64 {
	f, _ := n.Float64()
	return f
}

type provider struct {
	name string
	fn   func() (*IPInfo, error)
}

var providers = []provider{
	{"ip-api.com", func() (*IPInfo, error) {
		var r struct {
			Status      string  `json:"status"`
			Query       string  `json:"query"`
			Country     string  `json:"country"`
			CountryCode string  `json:"countryCode"`
			RegionName  string  `json:"regionName"`
			City        string  `json:"city"`
			Zip         string  `json:"zip"`
			Lat         float64 `json:"lat"`
			Lon         float64 `json:"lon"`
			Timezone    string  `json:"timezone"`
			ISP         string  `json:"isp"`
			Org         string  `json:"org"`
		}
		if err := getJSON("http://ip-api.com/json/?fields=status,query,country,countryCode,regionName,city,zip,lat,lon,timezone,isp,org", &r); err != nil {
			return nil, err
		}
		if r.Status != "success" {
			return nil, fmt.Errorf("status %s", r.Status)
		}
		isp := r.ISP
		if isp == "" {
			isp = r.Org
		}
		i := &IPInfo{IP: r.Query, Country: r.Country, CountryCode: r.CountryCode, City: r.City,
			Region: r.RegionName, Zip: r.Zip, Timezone: r.Timezone, ISP: isp}
		i.setGeo(r.Lat, r.Lon)
		return i, nil
	}},
	{"geojs.io", func() (*IPInfo, error) {
		var r struct {
			IP           string      `json:"ip"`
			Country      string      `json:"country"`
			CountryCode  string      `json:"country_code"`
			City         string      `json:"city"`
			Region       string      `json:"region"`
			Timezone     string      `json:"timezone"`
			Organization string      `json:"organization_name"`
			Lat          json.Number `json:"latitude"`
			Lon          json.Number `json:"longitude"`
		}
		if err := getJSON("https://get.geojs.io/v1/ip/geo.json", &r); err != nil {
			return nil, err
		}
		i := &IPInfo{IP: r.IP, Country: r.Country, CountryCode: r.CountryCode, City: r.City,
			Region: r.Region, Timezone: r.Timezone, ISP: r.Organization}
		i.setGeo(num(r.Lat), num(r.Lon))
		return i, nil
	}},
	{"ipwho.is", func() (*IPInfo, error) {
		var r struct {
			Success     bool    `json:"success"`
			IP          string  `json:"ip"`
			Country     string  `json:"country"`
			CountryCode string  `json:"country_code"`
			Region      string  `json:"region"`
			City        string  `json:"city"`
			Postal      string  `json:"postal"`
			Lat         float64 `json:"latitude"`
			Lon         float64 `json:"longitude"`
			Timezone    struct {
				ID string `json:"id"`
			} `json:"timezone"`
			Connection struct {
				ISP string `json:"isp"`
				Org string `json:"org"`
			} `json:"connection"`
		}
		if err := getJSON("https://ipwho.is/", &r); err != nil {
			return nil, err
		}
		if !r.Success {
			return nil, fmt.Errorf("unsuccessful")
		}
		isp := r.Connection.ISP
		if isp == "" {
			isp = r.Connection.Org
		}
		i := &IPInfo{IP: r.IP, Country: r.Country, CountryCode: r.CountryCode, City: r.City,
			Region: r.Region, Zip: r.Postal, Timezone: r.Timezone.ID, ISP: isp}
		i.setGeo(r.Lat, r.Lon)
		return i, nil
	}},
	{"ipapi.co", func() (*IPInfo, error) {
		var r struct {
			IP          string  `json:"ip"`
			CountryName string  `json:"country_name"`
			CountryCode string  `json:"country_code"`
			City        string  `json:"city"`
			Region      string  `json:"region"`
			Postal      string  `json:"postal"`
			Timezone    string  `json:"timezone"`
			Org         string  `json:"org"`
			Lat         float64 `json:"latitude"`
			Lon         float64 `json:"longitude"`
		}
		if err := getJSON("https://ipapi.co/json/", &r); err != nil {
			return nil, err
		}
		i := &IPInfo{IP: r.IP, Country: r.CountryName, CountryCode: r.CountryCode, City: r.City,
			Region: r.Region, Zip: r.Postal, Timezone: r.Timezone, ISP: r.Org}
		i.setGeo(r.Lat, r.Lon)
		return i, nil
	}},
	{"ipinfo.io", func() (*IPInfo, error) {
		var r struct {
			IP       string `json:"ip"`
			Country  string `json:"country"`
			City     string `json:"city"`
			Region   string `json:"region"`
			Postal   string `json:"postal"`
			Timezone string `json:"timezone"`
			Org      string `json:"org"`
			Loc      string `json:"loc"` // "lat,lon"
		}
		if err := getJSON("https://ipinfo.io/json", &r); err != nil {
			return nil, err
		}
		i := &IPInfo{IP: r.IP, Country: r.Country, CountryCode: r.Country, City: r.City,
			Region: r.Region, Zip: r.Postal, Timezone: r.Timezone, ISP: r.Org}
		if parts := strings.Split(r.Loc, ","); len(parts) == 2 {
			lat, e1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			lon, e2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if e1 == nil && e2 == nil {
				i.setGeo(lat, lon)
			}
		}
		return i, nil
	}},
	{"myip.com", func() (*IPInfo, error) { // без координат
		var r struct {
			IP      string `json:"ip"`
			Country string `json:"country"`
			CC      string `json:"cc"`
		}
		if err := getJSON("https://api.myip.com", &r); err != nil {
			return nil, err
		}
		return &IPInfo{IP: r.IP, Country: r.Country, CountryCode: r.CC}, nil
	}},
	{"ipify.org", func() (*IPInfo, error) { // только IP
		var r struct {
			IP string `json:"ip"`
		}
		if err := getJSON("https://api.ipify.org?format=json", &r); err != nil {
			return nil, err
		}
		return &IPInfo{IP: r.IP}, nil
	}},
}

type result struct {
	name string
	info *IPInfo
	err  error
}

// score: чем больше, тем богаче ответ. С координатами — максимум.
func score(i *IPInfo) int {
	switch {
	case i.HasGeo:
		return 3
	case i.Country != "" || i.CountryCode != "":
		return 2
	default:
		return 1
	}
}

// GetIPInfo вызывается из JS: window.go.main.App.GetIPInfo()
// Все сервисы опрашиваются параллельно. Побеждает первый ответ с координатами;
// если координат нет ни у кого — лучший из имеющихся.
func (a *App) GetIPInfo() (*IPInfo, error) {
	ch := make(chan result, len(providers))
	for _, p := range providers {
		go func(p provider) {
			start := time.Now()
			info, err := p.fn()
			if err == nil && (info == nil || info.IP == "") {
				err = fmt.Errorf("empty answer")
			}
			if err == nil {
				info.PingMs = int(time.Since(start).Milliseconds())
				info.Provider = p.name
			}
			ch <- result{p.name, info, err}
		}(p)
	}

	var best *IPInfo
	var errs []string
	for range providers {
		r := <-ch
		if r.err != nil {
			errs = append(errs, r.name+": "+r.err.Error())
			continue
		}
		if r.info.HasGeo {
			return r.info, nil
		}
		if best == nil || score(r.info) > score(best) {
			best = r.info
		}
	}
	if best != nil {
		return best, nil
	}
	return nil, fmt.Errorf("no service responded (%s)", strings.Join(errs, "; "))
}

// Place — найденная по тексту точка (для выбора начальной точки дуги).
type Place struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

// shortName: "Москва, Центральный федеральный округ, Россия" -> "Москва, Центральный федеральный округ, Россия"
// длинные адреса сокращаются до: первая часть, вторая часть, страна.
func shortName(full string) string {
	parts := strings.Split(full, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	if len(parts) > 3 {
		return parts[0] + ", " + parts[1] + ", " + parts[len(parts)-1]
	}
	return strings.Join(parts, ", ")
}

// GeocodeAddress вызывается из JS: window.go.main.App.GeocodeAddress("Москва, Тверская 1")
// Сначала OpenStreetMap Nominatim (понимает улицы и дома), затем Open-Meteo (только города).
func (a *App) GeocodeAddress(q string) (*Place, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil, fmt.Errorf("пустой адрес")
	}

	var nm []struct {
		Lat         string `json:"lat"`
		Lon         string `json:"lon"`
		DisplayName string `json:"display_name"`
	}
	u1 := "https://nominatim.openstreetmap.org/search?format=json&limit=1&accept-language=ru&q=" + url.QueryEscape(q)
	if err := getJSON(u1, &nm); err == nil && len(nm) > 0 {
		lat, e1 := strconv.ParseFloat(nm[0].Lat, 64)
		lon, e2 := strconv.ParseFloat(nm[0].Lon, 64)
		if e1 == nil && e2 == nil {
			return &Place{Name: shortName(nm[0].DisplayName), Lat: lat, Lon: lon}, nil
		}
	}

	var om struct {
		Results []struct {
			Name      string  `json:"name"`
			Country   string  `json:"country"`
			Admin1    string  `json:"admin1"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}
	u2 := "https://geocoding-api.open-meteo.com/v1/search?count=1&format=json&language=ru&name=" + url.QueryEscape(q)
	if err := getJSON(u2, &om); err == nil && len(om.Results) > 0 {
		r := om.Results[0]
		name := r.Name
		if r.Admin1 != "" && r.Admin1 != r.Name {
			name += ", " + r.Admin1
		}
		if r.Country != "" {
			name += ", " + r.Country
		}
		return &Place{Name: name, Lat: r.Latitude, Lon: r.Longitude}, nil
	}

	return nil, fmt.Errorf("адрес не найден (проверьте написание или интернет)")
}
