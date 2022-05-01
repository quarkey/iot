import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { flatMap } from "rxjs/operators";
import { environment } from "src/environments/environment";
import { Dashboard } from "../models/general";

@Injectable({
  providedIn: "root",
})
export class GeneralService {
  constructor(private http: HttpClient) {}

  public Dashboard() {
    return this.http.get<Dashboard>(`${environment.apiUrl}/api/dashboard`);
  }
  DownloadFile(data: Object, filename: string) {
    let csvData = this.convertToCSV(data);
    let blob = new Blob(["\ufeff" + csvData], {
      type: "text/csv;charset=utf-8;",
    });
    let dwldLink = document.createElement("a");
    let url = URL.createObjectURL(blob);
    let isSafariBrowser =
      navigator.userAgent.indexOf("Safari") != -1 &&
      navigator.userAgent.indexOf("Chrome") == -1;
    if (isSafariBrowser) {
      //if Safari open in new window to save file with random filename.
      dwldLink.setAttribute("target", "_blank");
    }
    dwldLink.setAttribute("href", url);
    dwldLink.setAttribute("download", filename);
    dwldLink.style.visibility = "hidden";
    document.body.appendChild(dwldLink);
    dwldLink.click();
    document.body.removeChild(dwldLink);
  }
  // ConvertToCSV is creating comma delimited string.
  private convertToCSV(objArray: Object) {
    let array = typeof objArray != "object" ? JSON.parse(objArray) : objArray;
    let str = "";
    array.forEach((element) => {
      str += element + "\r\n";
    });
    return str;
  }
}
