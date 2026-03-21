export namespace helper {
	
	export class SdkInfo {
	    path: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new SdkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.source = source["source"];
	    }
	}

}

export namespace models {
	
	export class AccelInfo {
	    status: string;
	    hypervisor: string;
	    details: string;
	
	    static createFrom(source: any = {}) {
	        return new AccelInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.hypervisor = source["hypervisor"];
	        this.details = source["details"];
	    }
	}
	export class AvdInfo {
	    name: string;
	    displayName: string;
	    path: string;
	    diskUsage: string;
	    running: boolean;
	    apiLevel: string;
	    androidVersion: string;
	    androidCodename: string;
	    abi: string;
	    ramSize: string;
	    resolution: string;
	    hasGooglePlay: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AvdInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.displayName = source["displayName"];
	        this.path = source["path"];
	        this.diskUsage = source["diskUsage"];
	        this.running = source["running"];
	        this.apiLevel = source["apiLevel"];
	        this.androidVersion = source["androidVersion"];
	        this.androidCodename = source["androidCodename"];
	        this.abi = source["abi"];
	        this.ramSize = source["ramSize"];
	        this.resolution = source["resolution"];
	        this.hasGooglePlay = source["hasGooglePlay"];
	    }
	}

}

export namespace services {
	
	export class GitHubRelease {
	    tag_name: string;
	    html_url: string;
	
	    static createFrom(source: any = {}) {
	        return new GitHubRelease(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag_name = source["tag_name"];
	        this.html_url = source["html_url"];
	    }
	}

}

