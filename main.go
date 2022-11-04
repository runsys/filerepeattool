// 目录单层化并哈希排重 project main.go
//author:iwlb@outlook.com
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"linbo.ga/toolfunc"
)

var nomove bool //no move to root directory

func GetDirSubFile(allbook []string, dir string, regex string) (mafs []string) {
	re := regexp.MustCompile(regex)
	for i := 0; i < len(allbook); i++ {
		if strings.HasPrefix(allbook[i], dir) {
			if strings.Index(allbook[i][len(dir):], "/") == -1 {
				if re.MatchString(allbook[i]) {
					mafs = append(mafs, allbook[i])
				}
			}
		}
	}
	if len(mafs) > 16 {
		log.Println("error:", len(mafs), dir, regex)
		panic("error")
	}
	return mafs
}

func DeleteRepeatFileByReserveLongName(txtfilepath string) {
	txtfilepath = toolfunc.StdFilePath(txtfilepath)
	txtdir := txtfilepath[:strings.LastIndex(txtfilepath, "/")+1]
	ctt, ctte := ioutil.ReadFile(txtfilepath)
	if ctte == nil {
		lis := strings.Split(string(ctt), "\n")
		hasn_files := make(map[string][]string, 0)
		for _, li := range lis {
			li = strings.Trim(li, " \r\n\t")
			if li == "" {
				continue
			}
			segs := strings.Split(li, "\t")
			hasn_files[segs[0]] = segs[1:]
		}
		allbook := []string{}
		toolfunc.GetDirAllFile("E:/book书籍/", ".*", &allbook)
		for _, li := range lis {
			li = strings.Trim(li, " \r\n\t")
			if li == "" {
				continue
			}
			segs := strings.Split(li, "\t")
			names := []string{}
			segs2 := [][]string{}
			var maxlenname string
			var maxlennameext string
			var maxlennamei int
			var is_nomove bool
			var existscnt int
			for _, seg := range segs {
				_, fie := os.Stat(seg)
				if fie == nil {
					existscnt += 1
				}
			}
			if existscnt <= 1 {
				continue
			}
			for si, seg := range segs {
				if si > 0 {
					filename := seg[strings.LastIndex(seg, "/")+1:]
					if strings.LastIndex(filename, ".") != -1 {
						filename = filename[:strings.LastIndex(filename, ".")]
					}
					names = append(names, filename)
					var ext string
					if strings.LastIndex(seg, ".") != -1 {
						ext = strings.ToLower(seg[strings.LastIndex(seg, "."):])
					}
					if len(filename) > len(maxlenname) {
						maxlenname = filename
						maxlennamei = si
						maxlennameext = ext
					}
					fidir := seg[:strings.LastIndex(seg, "/")+1]
					var noextpa2 string = seg
					if strings.LastIndex(seg, ".") != -1 && strings.LastIndex(seg, ".") > strings.LastIndex(seg, "/") {
						noextpa2 = seg[:strings.LastIndex(seg, ".")]
					}
					savenamefiles := []string{}
					noextpa3 := noextpa2
					noextpa3 = strings.Replace(noextpa3, "\\", "\\\\", -1)
					noextpa3 = strings.Replace(noextpa3, "[", "\\[", -1)
					noextpa3 = strings.Replace(noextpa3, "]", "\\]", -1)
					noextpa3 = strings.Replace(noextpa3, "*", "\\*", -1)
					noextpa3 = strings.Replace(noextpa3, "+", "\\+", -1)
					noextpa3 = strings.Replace(noextpa3, "{", "\\{", -1)
					noextpa3 = strings.Replace(noextpa3, "}", "\\}", -1)
					noextpa3 = strings.Replace(noextpa3, ".", "[.]", -1)
					noextpa3 = strings.Replace(noextpa3, "?", "\\?", -1)
					noextpa3 = strings.Replace(noextpa3, "(", "\\(", -1)
					noextpa3 = strings.Replace(noextpa3, ")", "\\)", -1)
					//fmt.Println("regex", noextpa2+"[.][a-zA-Z0-9]+")
					//log.Println("regex", noextpa2+"[.][a-zA-Z0-9]+")

					savenamefiles = GetDirSubFile(allbook, fidir, noextpa3+"[.][a-zA-Z0-9]+")

					_, fie := os.Stat(noextpa2 + "_files")
					if (ext == ".html" || ext == ".htm") && fie == nil {
						savenamefiles = append(savenamefiles, noextpa2+"_files")
					}
					segs2 = append(segs2, savenamefiles)
				} else {
					segs2 = append(segs2, []string{})
				}
			}
			groupcnt := 0
			for _, s2 := range segs2 {
				if len(s2) > 1 {
					groupcnt += 1
				}
			}
			if toolfunc.IsNumber(maxlenname) {
				is_nomove = true
			}
			if nomove {
				is_nomove = true
			}
			if is_nomove == true {
				if groupcnt == 0 {
					maxl := 0
					for j := 1; j < len(segs); j++ {
						if len(segs) > maxl {
							maxl = len(segs[j])
						}
					}
					skipcnt := 0
					for j := 1; j < len(segs); j++ {
						if len(segs) == maxl && skipcnt == 0 {
							skipcnt += 1
							continue
						} else {
							fmt.Println("remove file", segs[j])
							log.Println("remove file", segs[j])
							os.Remove(segs[j])
						}
					}
				} else if groupcnt == 1 {
					tarmap := make(map[string]int, 0)
					for _, s2 := range segs2 {
						if len(s2) > 1 {
							for j := 0; j < len(s2); j++ {
								s2fi, s2fie := os.Stat(s2[j])
								if s2fie != nil {
									continue
								}
								if s2fi.IsDir() == false {
									s2[j] = toolfunc.StdPath(s2[j])
									tarmap[s2[j]] = 1
								} else {
									if strings.HasSuffix(s2[j], "/") {
										s2[j] = s2[j][:len(s2[j])-1]
									}
									s2[j] = toolfunc.StdDir(s2[j])
									tarmap[s2[j]] = 1
								}
							}
						} else {
							// if len(s2[0])==1 {
							// 	os.Remove(s2[0])
							// }
						}
					}

					for _, s2 := range segs2 {
						if len(s2) > 1 {
						} else if len(s2) > 0 {
							if len(s2) != 1 {
								panic("err")
							}
							s2[0] = toolfunc.StdPath(s2[0])
							_, bintama := tarmap[s2[0]]
							if bintama == false {
								fmt.Println("remove file", s2[0])
								log.Println("remove file", s2[0])
								os.Remove(s2[0])
							}
						}
					}
				} else if groupcnt > 1 {
					maxcnt := 0
					for _, s2 := range segs2 {
						if len(s2) > 1 {
							if len(s2) > maxcnt {
								maxcnt = len(s2)
							}
						}
					}
					tarpa := make(map[string]int, 0)
					bsizezero := false
					for _, s2 := range segs2 {
						if len(s2) == maxcnt {
							for j := 0; j < len(s2); j++ {
								s2fi, s2fie := os.Stat(s2[j])
								if s2fie != nil {
									continue
								}
								if s2fi.Size() == 0 {
									bsizezero = true
								}
								if s2fi.IsDir() == false {
									s2[j] = toolfunc.StdPath(s2[j])
									tarpa[s2[j]] = 1
								} else {
									if strings.HasSuffix(s2[j], "/") {
										s2[j] = s2[j][:len(s2[j])-1]
									}
									s2[j] = toolfunc.StdPath(s2[j])
									tarpa[s2[j]] = 1
								}
							}
							break
						} else {
							// for j:=0;j<len(s2);j++ {
							// 	os.Remove(s2[j])
							// }
						}
					}
					if bsizezero == false {
						for _, s2 := range segs2 {
							for j := 0; j < len(s2); j++ {
								s2fi, s2fie := os.Stat(s2[j])
								if s2fie != nil {
									continue
								}
								if s2fi.IsDir() == false {
									s2[j] = toolfunc.StdPath(s2[j])
									_, beintarma := tarpa[s2[j]]
									if beintarma == false {
										fmt.Println("remove file", s2[j])
										log.Println("remove file", s2[j])
										os.Remove(s2[j])
									}
								} else {
									s2[j] = toolfunc.StdDir(s2[j])
									_, beintarma := tarpa[s2[j]]
									if beintarma == false {
										fmt.Println("remove dir all", s2[j])
										log.Println("remove dir all", s2[j])
										os.RemoveAll(s2[j])
									}
								}
							}
						}
					}
				}
			} else {
				if groupcnt == 0 {
					fmt.Println("move file", segs[maxlennamei], txtdir+maxlenname+maxlennameext)
					log.Println("move file", segs[maxlennamei], txtdir+maxlenname+maxlennameext)
					toolfunc.MoveFile(segs[maxlennamei], txtdir+maxlenname+maxlennameext)
					for j := 1; j < len(segs); j++ {
						if toolfunc.StdPath(segs[j]) == toolfunc.StdPath(txtdir+maxlenname+maxlennameext) {
							continue
						} else {
							fmt.Println("remove file", segs[j])
							log.Println("remove file", segs[j])
							os.Remove(segs[j])
						}
					}
				} else if groupcnt == 1 {
					tarmap := make(map[string]int, 0)
					for _, s2 := range segs2 {
						if len(s2) > 1 {
							for j := 0; j < len(s2); j++ {
								s2fi, s2fie := os.Stat(s2[j])
								if s2fie != nil {
									continue
								}
								if s2fi.IsDir() == false {
									movetoname := s2[j][strings.LastIndex(s2[j], "/")+1:]
									targetpath := txtdir + movetoname
									fmt.Println("move file", s2[j], targetpath)
									log.Println("move file", s2[j], targetpath)
									toolfunc.MoveFile(s2[j], targetpath)
									s2[j] = toolfunc.StdPath(s2[j])
									targetpath = toolfunc.StdPath(targetpath)
									tarmap[targetpath] = 1
									if s2[j] != targetpath {
										fmt.Println("remove file", s2[j])
										log.Println("remove file", s2[j])
										os.Remove(s2[j])
									}
								} else {
									if strings.HasSuffix(s2[j], "/") {
										s2[j] = s2[j][:len(s2[j])-1]
									}

									targetpath := txtdir + maxlenname + "_files"
									s2[j] = toolfunc.StdDir(s2[j])
									targetpath = toolfunc.StdDir(targetpath)
									fmt.Println("move dir", s2[j], targetpath)
									log.Println("move dir", s2[j], targetpath)
									toolfunc.MoveDir(s2[j], targetpath)
									tarmap[targetpath] = 1
									if s2[j] != targetpath {
										fmt.Println("remove dir all", s2[j])
										log.Println("remove dir all", s2[j])
										os.RemoveAll(s2[j])
									}
								}
							}
						} else {
							// if len(s2[0])==1 {
							// 	os.Remove(s2[0])
							// }
						}
					}

					for _, s2 := range segs2 {
						if len(s2) > 1 {
						} else if len(s2) > 0 {
							if len(s2) != 1 {
								panic("err")
							}
							s2[0] = toolfunc.StdPath(s2[0])
							_, bintama := tarmap[s2[0]]
							if bintama == false {
								fmt.Println("remove file", s2[0])
								log.Println("remove file", s2[0])
								os.Remove(s2[0])
							}
						}
					}
				} else if groupcnt > 1 {
					maxcnt := 0
					for _, s2 := range segs2 {
						if len(s2) > 1 {
							if len(s2) > maxcnt {
								maxcnt = len(s2)
							}
						}
					}
					tarpa := make(map[string]int, 0)
					bsizezero := false
					for _, s2 := range segs2 {
						if len(s2) == maxcnt {
							for j := 0; j < len(s2); j++ {
								s2fi, s2fie := os.Stat(s2[j])
								if s2fie != nil {
									continue
								}
								if s2fi.Size() == 0 {
									bsizezero = true
								}
								if s2fi.IsDir() == false {
									movetoname := s2[j][strings.LastIndex(s2[j], "/")+1:]
									targetpath := txtdir + movetoname
									fmt.Println("move file", s2[j], targetpath)
									log.Println("move file", s2[j], targetpath)
									toolfunc.MoveFile(s2[j], targetpath)
									s2[j] = toolfunc.StdPath(s2[j])
									targetpath = toolfunc.StdPath(targetpath)
									tarpa[targetpath] = 1
									if s2[j] != targetpath {
										fmt.Println("remove file", s2[j])
										log.Println("remove file", s2[j])
										os.Remove(s2[j])
									}
								} else {
									if strings.HasSuffix(s2[j], "/") {
										s2[j] = s2[j][:len(s2[j])-1]
									}
									targetpath := txtdir + maxlenname + "_files"
									s2[j] = toolfunc.StdPath(s2[j])
									targetpath = toolfunc.StdPath(targetpath)
									fmt.Println("move dir", s2[j], targetpath)
									log.Println("move dir", s2[j], targetpath)
									toolfunc.MoveDir(s2[j], targetpath)
									tarpa[targetpath] = 1
									if s2[j] != targetpath {
										fmt.Println("remove dir all", s2[j])
										log.Println("remove dir all", s2[j])
										os.RemoveAll(s2[j])
									}
								}
							}
							break
						} else {
							// for j:=0;j<len(s2);j++ {
							// 	os.Remove(s2[j])
							// }
						}
					}
					if bsizezero == false {
						for _, s2 := range segs2 {
							for j := 0; j < len(s2); j++ {
								s2fi, s2fie := os.Stat(s2[j])
								if s2fie != nil {
									continue
								}
								if s2fi.IsDir() == false {
									s2[j] = toolfunc.StdPath(s2[j])
									_, beintarma := tarpa[s2[j]]
									if beintarma == false {
										fmt.Println("remove file", s2[j])
										log.Println("remove file", s2[j])
										os.Remove(s2[j])
									}
								} else {
									s2[j] = toolfunc.StdDir(s2[j])
									_, beintarma := tarpa[s2[j]]
									if beintarma == false {
										fmt.Println("remove dir all", s2[j])
										log.Println("remove dir all", s2[j])
										os.RemoveAll(s2[j])
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func main() {
	fmt.Println("生成排重文件> [dir of check repeat]\n>执行排重> [目录查找到的重复文件及MD5.txt path] [nomove 是否移动到主目录以单层化]")
	log.Println("生成排重文件> [dir of check repeat]\n>执行排重> [目录查找到的重复文件及MD5.txt path] [nomove 是否移动到主目录以单层化]")
	if len(os.Args) == 1 {
		return
	}
	logf, _ := os.OpenFile("目录单层化并哈希排重.log", os.O_WRONLY|os.O_CREATE, 0666)
	logf.Seek(0, os.SEEK_END)
	defer logf.Close()
	log.SetOutput(logf)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.LstdFlags)

	if strings.HasSuffix(os.Args[1], "目录查找到的重复文件及MD5.txt") {
		nomove = false
		if len(os.Args) == 3 {
			if os.Args[2] == "nomove" {
				nomove = true
			}
		}
		DeleteRepeatFileByReserveLongName(os.Args[1])
		return
	}

	checkdir := toolfunc.StdDir(os.Args[1])
	_, fie := os.Stat(checkdir)
	if fie != nil {
		log.Println("目录不存在错误")
		return
	}

	allfile := []string{}
	toolfunc.GetDirAllFile(checkdir, "", &allfile)
	fmt.Println("获取列表完成")
	md5_paths := make(map[string][]string, 0)
	samename_path := make(map[string][]string, 0)
	net := time.Now().Add(30 * time.Second)
	var donecnt int
	donesize := 0
	for _, fil := range allfile {
		dfi, dfie := os.Stat(fil)
		if dfie != nil {
			continue
		}
		donesize += int(dfi.Size())
		filmd5 := toolfunc.FileMD5Hash(fil)
		pas, be := md5_paths[filmd5]
		if be {
			md5_paths[filmd5] = append(pas, fil)
		} else {
			md5_paths[filmd5] = []string{fil}
		}
		filname := fil[strings.LastIndex(fil, "/")+1:]
		if strings.LastIndex(filname, ".") != -1 {
			filname = filname[:strings.LastIndex(filname, ".")]
			olfs, be2 := samename_path[filname]
			if be2 {
				samename_path[filname] = append(olfs, fil)
			} else {
				samename_path[filname] = []string{fil}
			}
		}
		donecnt += 1
		if net.Before(time.Now()) {
			net = time.Now().Add(30 * time.Second)
			fmt.Println("total count", len(allfile), "done count", donecnt, "done size", donesize)
		}
	}

	var bts2 []byte
	var nocf []byte
	for k, v := range md5_paths {
		if len(v) > 1 {
			bts2 = append(bts2, []byte(k+"\t"+strings.Join(v, "\t")+"\n")...)
		} else {
			nocf = append(nocf, []byte(k+"\t"+v[0]+"\n")...)
		}
	}
	ioutil.WriteFile(checkdir+"目录查找到的重复文件及MD5.txt", bts2, 0666)
	ioutil.WriteFile(checkdir+"没有重复文件的MD5.txt", nocf, 0666)

	var samenamebs []byte
	for k, v := range md5_paths {
		if len(v) > 1 {
			samenamebs = append(samenamebs, []byte(k+"\t"+strings.Join(v, "\t")+"\n")...)
		}
	}
	ioutil.WriteFile(checkdir+"目录中名字相同的文件.txt", samenamebs, 0666)
}
