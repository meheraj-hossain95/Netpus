export namespace database {
	
	export class AppUsageStat {
	    AppName: string;
	    TotalUpload: number;
	    TotalDownload: number;
	    LastSeen: number;
	
	    static createFrom(source: any = {}) {
	        return new AppUsageStat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppName = source["AppName"];
	        this.TotalUpload = source["TotalUpload"];
	        this.TotalDownload = source["TotalDownload"];
	        this.LastSeen = source["LastSeen"];
	    }
	}
	export class DailySummary {
	    ID: number;
	    Date: string;
	    TotalUpload: number;
	    TotalDownload: number;
	
	    static createFrom(source: any = {}) {
	        return new DailySummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Date = source["Date"];
	        this.TotalUpload = source["TotalUpload"];
	        this.TotalDownload = source["TotalDownload"];
	    }
	}

}

export namespace monitor {
	
	export class MonitorStatus {
	    running: boolean;
	    paused: boolean;
	    updateInterval: number;
	    // Go type: time
	    lastUpdate: any;
	
	    static createFrom(source: any = {}) {
	        return new MonitorStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.paused = source["paused"];
	        this.updateInterval = source["updateInterval"];
	        this.lastUpdate = this.convertValues(source["lastUpdate"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace utils {
	
	export class Config {
	    autoStart: boolean;
	    theme: string;
	    dataRetention: number;
	    networkInterface: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.autoStart = source["autoStart"];
	        this.theme = source["theme"];
	        this.dataRetention = source["dataRetention"];
	        this.networkInterface = source["networkInterface"];
	    }
	}

}

