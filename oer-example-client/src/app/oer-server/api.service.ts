import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Channel, ProgramEntry, StatusResponse} from "./entities";

@Injectable()
export class ApiService {

  private apiEndpoint: string = '/v1';

  constructor(private http: HttpClient) { }

  public statusResponse(): Observable<StatusResponse> {
    return this.http.get<StatusResponse>(this.apiEndpoint + "/status");
  }

  public channels(): Observable<Channel[]> {
    return this.http.get<Channel[]>(this.apiEndpoint + "/channels");
  }

  public dailyProgram(): Observable<ProgramEntry[]> {
    return this.http.get<ProgramEntry[]>(this.apiEndpoint + "/program/daily");
  }

  public program(from: Date, to: Date): Observable<ProgramEntry[]> {
    return this.http.get<ProgramEntry[]>(this.apiEndpoint + "/program?from=" + from.toISOString() + "&to=" + to.toISOString());
  }

  public entry(clickedEntryId: number) : Observable<ProgramEntry> {
    return this.http.get<ProgramEntry>(this.apiEndpoint + '/program/entry/' + clickedEntryId);
  }
}
