export namespace main {
	
	export class ClientInfo {
	    id: string;
	    hostname: string;
	    ip: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ClientInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.hostname = source["hostname"];
	        this.ip = source["ip"];
	        this.status = source["status"];
	    }
	}
	export class SystemInfo {
	    hostname: string;
	    user: string;
	    os: string;
	    cpu: number;
	    memMB: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.user = source["user"];
	        this.os = source["os"];
	        this.cpu = source["cpu"];
	        this.memMB = source["memMB"];
	    }
	}

}

