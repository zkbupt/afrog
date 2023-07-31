package config

import (
	"fmt"
	"github.com/zkbupt/afrog/pkg/log"
	"github.com/zkbupt/afrog/pkg/output"
	"github.com/zkbupt/afrog/pkg/poc"
	"github.com/zkbupt/afrog/pkg/utils"
	"github.com/zkbupt/afrog/pocs"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/zan8in/goflags"
	"github.com/zan8in/gologger"
	fileutil "github.com/zan8in/pins/file"
	sliceutil "github.com/zan8in/pins/slice"
)

var (
	ReverseCeyeApiKey string
	ReverseCeyeDomain string
	ReverseJndi       string
	ReverseLdapPort   string
	ReverseApiPort    string

	ReverseCeyeLive bool
	ReverseJndiLive bool
)

type Options struct {
	// afrog-config.yaml configuration file
	Config *Config

	// Pocs Directory
	PocsDirectory utils.StringSlice

	Targets sliceutil.SafeSlice

	// target URLs/hosts to scan
	Target goflags.StringSlice

	// list of target URLs/hosts to scan (one per line)
	TargetsFile string

	// PoC file or directory to scan
	PocFile string

	// Append PoC file or directory to scan
	AppendPoc goflags.StringSlice

	// show afrog-pocs list
	PocList bool

	// show a afrog-pocs detail
	PocDetail string

	ExcludePocs     goflags.StringSlice
	ExcludePocsFile string

	// file to write output to (optional), support format: html
	Output string

	// file to write output to (optional), support format: json
	Json string

	// file to write output to (optional), support format: json
	JsonAll string

	// search PoC by keyword , eg: -s tomcat
	Search string

	SearchKeywords []string

	// no progress if silent is true
	Silent bool

	// pocs to run based on severity. Possible values: info, low, medium, high, critical
	Severity string

	SeverityKeywords []string

	// update afrog-pocs
	UpdatePocs bool

	// update afrog version
	Update bool

	// Disable update check
	DisableUpdateCheck bool

	MonitorTargets bool

	// POC Execution Duration Tracker
	PocExecutionDurationMonitor bool

	// Single Vulnerability Stopper
	VulnerabilityScannerBreakpoint bool

	// Scan count num(targets * allpocs)
	Count int

	// Current Scan count num
	CurrentCount uint32

	// Thread lock
	OptLock sync.Mutex

	// Callback scan result
	// OnResult OnResult

	// maximum number of requests to send per second (default 150)
	RateLimit int

	// maximum number of afrog-pocs to be executed in parallel (default 25)
	Concurrency int

	// Smart Control Concurrency
	Smart bool

	// number of times to retry a failed request (default 1)
	Retries int

	//
	MaxHostError int

	// time to wait in seconds before timeout (default 10)
	Timeout int

	// http/socks5 proxy to use
	Proxy string

	MaxRespBodySize int

	// afrog process count (target total × pocs total)
	ProcessTotal uint32

	DisableOutputHtml bool

	OJ *output.OutputJson
}

func NewOptions() (*Options, error) {

	options := &Options{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`afrog`)

	flagSet.CreateGroup("target", "Target",
		flagSet.StringSliceVarP(&options.Target, "target", "t", nil, "target URLs/hosts to scan (comma separated)", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.TargetsFile, "target-file", "T", "", "list of target URLs/hosts to scan (one per line)"),
	)

	flagSet.CreateGroup("pocs", "PoCs",
		flagSet.StringVarP(&options.PocFile, "poc-file", "P", "", "PoC file or directory to scan"),
		flagSet.StringSliceVarP(&options.AppendPoc, "append-poc", "ap", nil, "append PoC file or directory to scan (comma separated)", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.PocDetail, "poc-detail", "pd", "", "show a afrog-pocs detail"),
		flagSet.BoolVarP(&options.PocList, "poc-list", "pl", false, "show afrog-pocs list"),
		flagSet.StringSliceVarP(&options.ExcludePocs, "exclude-pocs", "ep", nil, "pocs to exclude from the scan (comma-separated)", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.ExcludePocsFile, "exclude-pocs-file", "epf", "", "list of pocs to exclude from scan (file)"),
	)

	flagSet.CreateGroup("output", "Output",
		flagSet.StringVarP(&options.Output, "output", "o", "", "write to the HTML file, including all vulnerability results"),
		flagSet.StringVarP(&options.Json, "json", "j", "", "write to the JSON file, but it will not include the request and response content"),
		flagSet.StringVarP(&options.JsonAll, "json-all", "ja", "", "write to the JSON file, including all vulnerability results"),
		flagSet.BoolVarP(&options.DisableOutputHtml, "disable-output-html", "doh", false, "disable the automatic generation of HTML reports (higher priority than the -o command)"),
	)

	flagSet.CreateGroup("filter", "Filter",
		flagSet.StringVarP(&options.Search, "search", "s", "", "search PoC by keyword , eg: -s tomcat,phpinfo"),
		flagSet.StringVarP(&options.Severity, "severity", "S", "", "pocs to run based on severity. support: info, low, medium, high, critical, unknown"),
	)

	flagSet.CreateGroup("rate-limit", "Rate-Limit",
		flagSet.IntVarP(&options.RateLimit, "rate-limit", "rl", 150, "maximum number of requests to send per second"),
		flagSet.IntVarP(&options.Concurrency, "concurrency", "c", 25, "maximum number of afrog-pocs to be executed in parallel"),
		flagSet.BoolVar(&options.Smart, "smart", false, "intelligent adjustment of concurrency based on changes in the total number of assets being scanned"),
	)

	flagSet.CreateGroup("optimization", "Optimization",
		flagSet.IntVar(&options.Retries, "retries", 1, "number of times to retry a failed request"),
		flagSet.IntVar(&options.Timeout, "timeout", 10, "time to wait in seconds before timeout"),
		flagSet.BoolVar(&options.MonitorTargets, "mt", false, "enable the monitor-target feature during scanning"),
		flagSet.IntVar(&options.MaxHostError, "mhe", 3, "max errors for a host before skipping from scan"),
		flagSet.IntVar(&options.MaxRespBodySize, "mrbs", 2, "max of http response body size"),
		flagSet.BoolVar(&options.Silent, "silent", false, "only results only"),
		flagSet.BoolVar(&options.PocExecutionDurationMonitor, "pedm", false, "This monitor tracks and records the execution time of each POC to identify the POC with the longest execution time."),
		flagSet.BoolVar(&options.VulnerabilityScannerBreakpoint, "vsb", false, "Once a vulnerability is detected, the scanning program will immediately halt the scan and report the identified vulnerability."),
	)

	flagSet.CreateGroup("update", "Update",
		flagSet.BoolVarP(&options.Update, "update", "un", false, "update afrog engine to the latest released version"),
		// flagSet.BoolVarP(&options.UpdatePocs, "update-pocs", "up", false, "update afrog-pocs to the latest released version"),
		flagSet.BoolVarP(&options.DisableUpdateCheck, "disable-update-check", "duc", false, "disable automatic afrog-pocs update check"),
	)

	flagSet.CreateGroup("proxy", "Proxy",
		flagSet.StringVar(&options.Proxy, "proxy", "", "list of http/socks5 proxy to use (comma separated or file input)"),
	)

	_ = flagSet.Parse()

	if err := options.verifyOptions(); err != nil {
		return options, err
	}

	return options, nil
}

func (opt *Options) verifyOptions() error {

	config, err := NewConfig()
	if err != nil {
		return err
	}
	opt.Config = config

	// init append poc
	if len(opt.AppendPoc) > 0 {
		poc.InitLocalAppendList(opt.AppendPoc)
	}

	// init test poc
	if len(opt.PocFile) > 0 {
		poc.InitLocalTestList([]string{opt.PocFile})

	}

	// initialized embed poc、local poc and append poc
	if len(pocs.EmbedFileList) == 0 && len(poc.LocalFileList) == 0 && len(poc.LocalAppendList) == 0 && len(poc.LocalTestList) == 0 {
		return fmt.Errorf("PoCs is not empty")
	}

	if opt.PocList {
		err := opt.PrintPocList()
		if err != nil {
			gologger.Error().Msg(err.Error())
		}
		os.Exit(0)
	}

	if len(opt.PocDetail) > 0 {
		opt.ReadPocDetail()
		os.Exit(0)
	}

	if opt.Update {
		err := updateEngine()
		if err != nil {
			gologger.Error().Msg(err.Error())
		}
		os.Exit(0)
	}

	au, err := NewAfrogUpdate(true)
	if err != nil {
		return err
	}

	if !opt.DisableUpdateCheck {
		info, _ := au.AfrogUpdatePocs()
		if len(info) > 0 && opt.UpdatePocs {
			gologger.Info().Msg(info)
		}
	}

	if len(opt.Target) == 0 && len(opt.TargetsFile) == 0 {
		return fmt.Errorf("either `target` or `target-file` must be set")
	}

	ShowBanner(au)

	if (len(opt.Config.Reverse.Ceye.Domain) == 0 && len(opt.Config.Reverse.Ceye.ApiKey) == 0) ||
		(len(opt.Config.Reverse.Jndi.JndiAddress) == 0 && len(opt.Config.Reverse.Jndi.LdapPort) == 0 && len(opt.Config.Reverse.Jndi.ApiPort) == 0) {
		homeDir, _ := os.UserHomeDir()
		configDir := strings.ReplaceAll(homeDir+"/.config/afrog/afrog-config.yaml", "\\", "/")
		gologger.Info().Msg("The reverse connection platform is not configured, which may affect the validation of certain RCE PoCs")
		gologger.Info().Msgf("Go to [%s] to configure the reverse connection platform\n", configDir)
		gologger.Info().Msg("Tutorial: https://github.com/zan8in/afrog/wiki/Configuration")
		gologger.Print().Msg("")
	}

	ReverseCeyeApiKey = opt.Config.Reverse.Ceye.ApiKey
	ReverseCeyeDomain = opt.Config.Reverse.Ceye.Domain

	ReverseJndi = opt.Config.Reverse.Jndi.JndiAddress
	ReverseLdapPort = opt.Config.Reverse.Jndi.LdapPort
	ReverseApiPort = opt.Config.Reverse.Jndi.ApiPort

	return nil
}

func (o *Options) SetSearchKeyword() bool {
	if len(o.Search) > 0 {
		arr := strings.Split(o.Search, ",")
		if len(arr) > 0 {
			for _, v := range arr {
				o.SearchKeywords = append(o.SearchKeywords, strings.TrimSpace(v))
			}
			return true
		}
	}
	return false
}

func (o *Options) CheckPocKeywords(id, name string) bool {
	if len(o.SearchKeywords) > 0 {
		for _, v := range o.SearchKeywords {
			v = strings.ToLower(v)
			if strings.Contains(strings.ToLower(id), v) || strings.Contains(strings.ToLower(name), v) {
				return true
			}
		}
	}
	return false
}

func (o *Options) SetSeverityKeyword() bool {
	if len(o.Severity) > 0 {
		arr := strings.Split(o.Severity, ",")
		if len(arr) > 0 {
			for _, v := range arr {
				o.SeverityKeywords = append(o.SeverityKeywords, strings.TrimSpace(v))
			}
			return true
		}
	}
	return false
}

func (o *Options) CheckPocSeverityKeywords(severity string) bool {
	if len(o.SeverityKeywords) > 0 {
		for _, v := range o.SeverityKeywords {
			if strings.EqualFold(severity, v) {
				return true
			}
		}
	}
	return false
}

func (o *Options) FilterPocSeveritySearch(pocId, pocInfoName, severity string) bool {
	var isShow bool
	if len(o.Search) > 0 && o.SetSearchKeyword() && len(o.Severity) > 0 && o.SetSeverityKeyword() {
		if o.CheckPocKeywords(pocId, pocInfoName) && o.CheckPocSeverityKeywords(severity) {
			isShow = true
		}
	} else if len(o.Severity) > 0 && o.SetSeverityKeyword() {
		if o.CheckPocSeverityKeywords(severity) {
			isShow = true
		}
	} else if len(o.Search) > 0 && o.SetSearchKeyword() {
		if o.CheckPocKeywords(pocId, pocInfoName) {
			isShow = true
		}
	} else {
		isShow = true
	}
	return isShow
}

func (o *Options) PrintPocList() error {

	var number = 1

	if len(pocs.EmbedFileList) > 0 {
		gologger.Print().Msg("---------- Embed PoCs -----------------")
		for _, v := range pocs.EmbedFileList {
			if poc, err := pocs.EmbedReadPocByPath(v); err == nil {
				if o.FilterPocSeveritySearch(poc.Id, poc.Info.Name, poc.Info.Severity) {
					gologger.Print().Msgf("%s [%s][%s][%s] author:%s\n",
						log.LogColor.Time(number),
						log.LogColor.Title(poc.Id),
						log.LogColor.Green(poc.Info.Name),
						log.LogColor.GetColor(poc.Info.Severity, poc.Info.Severity), poc.Info.Author)
					number++
				}
			}
		}
	}

	// init LocalPocsDirectory
	if len(poc.LocalFileList) > 0 {
		gologger.Print().Msg("---------- Local afrog-pocs -----------------")
		for _, v := range poc.LocalFileList {
			if poc, err := poc.LocalReadPocByPath(v); err == nil {
				if o.FilterPocSeveritySearch(poc.Id, poc.Info.Name, poc.Info.Severity) {
					gologger.Print().Msgf("%s [%s][%s][%s] author:%s\n",
						log.LogColor.Time(number),
						log.LogColor.Title(poc.Id),
						log.LogColor.Green(poc.Info.Name),
						log.LogColor.GetColor(poc.Info.Severity, poc.Info.Severity), poc.Info.Author)
					number++
				}
			}
		}
	}

	// append pocs
	if len(poc.LocalAppendList) > 0 {
		gologger.Print().Msg("---------- Local append-pocs -----------------")
		for _, v := range poc.LocalAppendList {
			if poc, err := poc.LocalReadPocByPath(v); err == nil {
				if o.FilterPocSeveritySearch(poc.Id, poc.Info.Name, poc.Info.Severity) {
					gologger.Print().Msgf("%s [%s][%s][%s] author:%s\n",
						log.LogColor.Time(number),
						log.LogColor.Title(poc.Id),
						log.LogColor.Green(poc.Info.Name),
						log.LogColor.GetColor(poc.Info.Severity, poc.Info.Severity), poc.Info.Author)
					number++
				}
			}
		}
	}

	gologger.Print().Msgf("--------------------------------\r\nTotal: %d\n", number-1)

	return nil
}

func (o *Options) ReadPocDetail() {
	if content, err := pocs.EmbedReadContentByName(o.PocDetail); err == nil && len(content) > 0 {
		gologger.Print().Msgf("%s\n", string(content))
		return
	}
	if content, err := poc.LocalReadContentByName(o.PocDetail); err == nil && len(content) > 0 {
		gologger.Print().Msgf("%s\n", string(content))
		return
	}
}

func (o *Options) CreatePocList() []poc.Poc {
	var pocSlice []poc.Poc

	if len(o.PocFile) > 0 && len(poc.LocalTestList) > 0 {
		for _, pocYaml := range poc.LocalTestList {
			if p, err := poc.LocalReadPocByPath(pocYaml); err == nil {
				pocSlice = append(pocSlice, p)
			}
		}
		return pocSlice
	}

	for _, pocYaml := range poc.LocalAppendList {
		if p, err := poc.LocalReadPocByPath(pocYaml); err == nil {
			pocSlice = append(pocSlice, p)
		}
	}

	for _, pocYaml := range poc.LocalFileList {
		if p, err := poc.LocalReadPocByPath(pocYaml); err == nil {
			pocSlice = append(pocSlice, p)

		}
	}

	for _, pocEmbedYaml := range pocs.EmbedFileList {
		if p, err := pocs.EmbedReadPocByPath(pocEmbedYaml); err == nil {
			pocSlice = append(pocSlice, p)
		}
	}

	newPocSlice := []poc.Poc{}
	for _, poc := range pocSlice {
		if o.FilterPocSeveritySearch(poc.Id, poc.Info.Name, poc.Info.Severity) {
			newPocSlice = append(newPocSlice, poc)
		}
	}

	latestPocSlice := []poc.Poc{}
	order := []string{"info", "low", "medium", "high", "critical"}
	for _, o := range order {
		for _, s := range newPocSlice {
			if o == strings.ToLower(s.Info.Severity) {
				latestPocSlice = append(latestPocSlice, s)
			}
		}
	}

	// exclude pocs
	excludePocs, _ := o.parseExcludePocs()
	finalPocSlice := []poc.Poc{}
	for _, poc := range latestPocSlice {
		if !isExcludePoc(poc, excludePocs) {
			finalPocSlice = append(finalPocSlice, poc)
		}
	}

	return finalPocSlice
}

func (o *Options) SmartControl() {
	numCPU := runtime.NumCPU()
	targetLen := o.Targets.Len()

	if o.Concurrency == 25 && targetLen <= 10 {
		o.Concurrency = 10
	} else if o.Concurrency == 25 && targetLen >= 1000 {
		o.Concurrency = numCPU * 30
	} else if o.Concurrency == 25 && targetLen >= 500 {
		o.Concurrency = numCPU * 20
	} else if o.Concurrency == 25 && targetLen >= 100 {
		o.Concurrency = numCPU * 10
	}
}

func (o *Options) parseExcludePocs() ([]string, error) {
	var excludePocs []string
	if len(o.ExcludePocs) > 0 {
		excludePocs = append(excludePocs, o.ExcludePocs...)
	}

	if len(o.ExcludePocsFile) > 0 {
		cdata, err := fileutil.ReadFile(o.ExcludePocsFile)
		if err != nil {
			if len(excludePocs) > 0 {
				return excludePocs, nil
			} else {
				return excludePocs, err
			}
		}
		for poc := range cdata {
			excludePocs = append(excludePocs, poc)
		}
	}
	return excludePocs, nil
}

func isExcludePoc(poc poc.Poc, excludePocs []string) bool {
	if len(excludePocs) == 0 {
		return false
	}
	for _, ep := range excludePocs {
		v := strings.ToLower(ep)
		if strings.Contains(strings.ToLower(poc.Id), v) || strings.Contains(strings.ToLower(poc.Info.Name), v) {
			return true
		}
	}
	return false
}
