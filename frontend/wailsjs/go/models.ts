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

