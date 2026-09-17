# 공개 함수와 메서드

GoPD를 다른 Go 패키지에서 사용할 때 호출할 수 있는 API 목록입니다. 공개 함수는 **9개**, 공개 메서드는 결과 타입의 것입니다. 타입·필드·상수 전체는 `go doc -all .`로 확인합니다.

모듈 경로는 `github.com/MyungSub0519/gopd`, 패키지 이름은 `gopd`입니다. 구문 분석·xref·객체 로딩·콘텐츠 해석 엔진은 `internal/` 아래에 있어 외부에서 import할 수 없고, 이 문서의 면만이 호출 계약입니다.

```go
import "github.com/MyungSub0519/gopd"
```

## 세 계약

| 진입점 | 반환 | 하는 일 | 보관 |
| --- | --- | --- | --- |
| `ParsePDF(path)` | `*PDF` | 텍스트·그래픽만 해석해 페이지별로 묶음 반환 | 없음 |
| `Open(path)`, `Read(r, size)` | `*DetailedPDF` | 전체 해석: 글리프·그래픽·이미지·폰트·주석·연산·출처 | 상세 결과 전체 |
| `Extract(path, options)`, `ExtractReader(r, size, options)` | `*Extraction` | `ExtractOptions` 플래그가 선택한 것만 생성 | Provenance 켤 때만 스냅샷 |

세 계약은 같은 해석기를 공유하지만 서로 다른 결과 타입을 가지며, 슬림 계약은 헤비 엔진 상태를 유지하지 않습니다.

## 함수 — 9개

| 시그니처 | 용도 |
| --- | --- |
| `func ParsePDF(path string) (*PDF, error)` | 기본 결과. 부분 결과와 오류가 함께 올 수 있으므로 반드시 err 확인 |
| `func Open(path string) (*DetailedPDF, error)` | 파일 경로로 상세 분석. 파일을 닫음 |
| `func Read(r io.ReaderAt, size int64) (*DetailedPDF, error)` | ReaderAt 스냅샷으로 상세 분석. r을 닫지 않음 |
| `func Extract(path string, options ExtractOptions) (*Extraction, error)` | 선택적 추출. 파일을 닫음 |
| `func ExtractReader(r io.ReaderAt, size int64, options ExtractOptions) (*Extraction, error)` | ReaderAt 형태의 Extract |
| `func Int(object Object) (int64, error)` | PDF Integer를 int64로 변환 |
| `func Number(object Object) (float64, error)` | Integer/Real을 유한한 float64로 변환 |
| `func IsStream(object Object) bool` | 객체 값이 Stream인지 확인 |
| `func IdentityMatrix() Matrix` | 항등 행렬 `[1 0 0 1 0 0]` |

## 결과 수명과 오류 처리

- 파일 입력은 메모리 스냅샷으로 보관합니다. 반환 결과에 `Close()`가 필요 없습니다.
- 파싱·콘텐츠 해석 실패 시 부분 결과와 오류가 함께 반환될 수 있으므로, 결과가 nil이 아니어도 반드시 오류를 확인합니다.
- 미지원 효과는 오류 없이 진단으로 남을 수도 있습니다. 기본 결과의 `PDF.Diagnostics`, 상세 결과의 `DetailedPDF.Diagnostics`, 추출 결과의 `Extraction.Diagnostics`를 확인하세요.
- 반환 구조체·슬라이스·공유 리소스는 읽기 전용으로 취급합니다.
- `Dictionary.Get` 키 누락은 `gopd.ErrMissingKey`를 감싼 오류입니다.

관련 문서: [선택적 추출](selective-extraction.md) · [기본 PDF API](basic-pdf.md) · [리소스 제한](resource-limits.md)
