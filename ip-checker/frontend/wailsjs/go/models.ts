export namespace main {
	
	export class IPInfo {
	    ip: string;
	    country: string;
	    countryCode: string;
	    city: string;
	    region: string;
	    zip: string;
	    timezone: string;
	    isp: string;
	    lat: number;
	    lon: number;
	    hasGeo: boolean;
	    pingMs: number;
	    provider: string;
	
	    static createFrom(source: any = {}) {
	        return new IPInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.country = source["country"];
	        this.countryCode = source["countryCode"];
	        this.city = source["city"];
	        this.region = source["region"];
	        this.zip = source["zip"];
	        this.timezone = source["timezone"];
	        this.isp = source["isp"];
	        this.lat = source["lat"];
	        this.lon = source["lon"];
	        this.hasGeo = source["hasGeo"];
	        this.pingMs = source["pingMs"];
	        this.provider = source["provider"];
	    }
	}
	export class Place {
	    name: string;
	    lat: number;
	    lon: number;
	
	    static createFrom(source: any = {}) {
	        return new Place(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.lat = source["lat"];
	        this.lon = source["lon"];
	    }
	}

}

