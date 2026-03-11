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
	
	export class AvdInfo {
	    name: string;
	    path: string;
	    diskUsage: string;
	
	    static createFrom(source: any = {}) {
	        return new AvdInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.diskUsage = source["diskUsage"];
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

