//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package mongodb

const clusterCertYaml = `
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZwekNDQTQrZ0F3SUJBZ0lVYzRvNWZCbVNINnowVzRDb3BvZEU4a3ZWS0pNd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1lqRUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdNQ2tOaGJHbG1iM0p1YVdFeEZEQVNCZ05WQkFjTQpDMHh2Y3lCQmJtZGxiR1Z6TVJFd0R3WURWUVFLREFoU2IyOTBJRWx1WXpFVk1CTUdBMVVFQXd3TWJXOXVaMjlrCllpNXZibXg1TUNBWERUSXdNREl4TnpBNU5Ea3hOVm9ZRHpJeE1qQXdNVEkwTURrME9URTFXakJpTVFzd0NRWUQKVlFRR0V3SlZVekVUTUJFR0ExVUVDQXdLUTJGc2FXWnZjbTVwWVRFVU1CSUdBMVVFQnd3TFRHOXpJRUZ1WjJWcwpaWE14RVRBUEJnTlZCQW9NQ0ZKdmIzUWdTVzVqTVJVd0V3WURWUVFEREF4dGIyNW5iMlJpTG05dWJIa3dnZ0lpCk1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQ0R3QXdnZ0lLQW9JQ0FRQ3JaL1lmMlFrZFNzdmQ2OUJXVDk1WUQ3VGwKU1RWdGhnQjlwbzFpVG5FMnd6VXpQSmdsNWZTMGpqdWYwZFJVMHloSk5wZkdrYUxwTmJWVEFIdkNTeUFRMmtJeQpvZTFkeWdQS2ZEdU9aQTZnTWF5K1p3ajNIZlB4SmJxUXNmWkJWMUkwblRXZGxsYlBhdGQ3ZGN0RE1uK3g2QStxCkJNZnhxT1RiZm9FM3YwTmhNQXIrdU5YaTAvd3F2UzBEVUxVR1VxWW1lT3drL3pYSGRNdS9rQ21MbEVzQ2JnYy8Ka3NZbEtvM0lUMnV3RlNtQ0xqa3U5VnVkV0lFa2hwWkRrZUlDUlNRcm1ISGpFVjRMNzZ6bmxYZ3pzdlozanlXNApQSkwzMnE2NWhpazVJTkpNNDdmL0pxVnBJZ2gxeGQxWnJReTZISE5rWU9DVzI3R3h6cUg4YjZjempkZDBnZ0xBCnlFRDM0UmFJazNBYkVldUtyeTR2ZGVxVHFzeU9HdnFCRUdYYlVibEtPMENhTXFPMi8zSW9aVlNFM05qOVhJRHEKQUVpSnplc3g1b3FaWW4rQzlOQm1JZXNMU1JBdTZidk5KZmdleVJkVWgvcnliYzZmMWJLUkVJRFhCcXZCQXg5agpNUXE5aVlERGhvSVhyNnA4anVCdFdPS1YxVVZObVV3WGtabGh3WVhkMGhhVSt4VnR1SzhKbThlNThpQjBXcytzCjBrV3ZMNVY4aUlSYVdEZ3pmdHYwNFgvL0VzUjAzcFd6OGhTLzV6dEJyMU1wTHl2NHFjNGY5N1ludDZmc3h2cXcKeGx0cmtXSUNsdmRobk9Sa0xDa0ZnY3Fna3JydDZoSGFzS2ZEcGprSjExQXVOS1drbU9wU2lQcTJqMCtOYXJWYQpSYmtrVmcvVkhXSlB5UXNwNVFJREFRQUJvMU13VVRBZEJnTlZIUTRFRmdRVWhsTkV5Y0V0K2pITXhqbHV0OHplCkNvbzREcU13SHdZRFZSMGpCQmd3Rm9BVWhsTkV5Y0V0K2pITXhqbHV0OHplQ29vNERxTXdEd1lEVlIwVEFRSC8KQkFVd0F3RUIvekFOQmdrcWhraUc5dzBCQVFzRkFBT0NBZ0VBSEZ6Vm51RkZ6ajRHZ0RnOXUyWHBROTFtNWtVagpVWitDUk1IMXhKMDNsMkxNeEVwNmY1M2ZrTnRvSWVLYUlZMVN0ZzE5TVZEYmM3OThFSStuQnRvOGowVTBRZFkwCjg5RFp3Qzg1dURPdzZINGRydFl1dUJ2U1o1QzdVZi9Bb1BKazZZN0E5V3NtZlFYNHgyVkd5YUx6SDNtQlJiOXQKaHo3ZkZBQzJCOUtDNzJMY3hObm9kK2hmZXplMTcvdVZGT2JYbjFjbW15bnNDd1pqWVJVRHVER09mWDZKSjhKMAprVUJCTFI0cno3ZXhOanhRajZSMWpVS3NycEhlSmtRSVF4MGhKc2VNRFNZM2FHczFGYk90SjAvWkdDYjVqSENPClFZZlZBejdhVFltS0U3QWJEVEhmcFd1cEEyOFFqU3Frd1RTSGkxYUc4VzB2Zmd4OVpXc1VwZjZ5dlZSZkdzcHcKUmMvRmtHVEJ0cXJDT0ZKN2JTNzNyaC9UYXpsendGTEhPaHVKaUczc1A4NDJ3cWJXSWNvYm01TU1qSElYeDhzWQo5L0w5eHd1REczMThuVFllNVdidWJqYVBUUGpkelJRYzA0TUJGY0YrTGFsM0tTbjEzRmlLZGtuQjh4RnI0TUlxClJQZTR5LzhlNGJZZ3dhMjFic3ZhZitJVjRoVDcvL0xNTEl3RG5hZmZNU1BySXprTGQzdG9xOTdvOHhGcmxBV0UKUmdjWlZpd0IxVVpBSlcrQ1RHdE91RDdJMFRpZ0RkSUlzY05nVzV5REQ1d29sSGRudlorYVU2d1dnQnJVNzJ4cAp4U2tnVyt5RmFybFBkbmFvSU1yUXRkcWJuSXFkYVJwdVY0N3QrT2lwN0hRUEJVNVN3d25aNGp5N21ZcnJvMWt0CjZzRk44ZnZoVnNBZit1Zz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUpRd0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQ1Mwd2dna3BBZ0VBQW9JQ0FRQ3JaL1lmMlFrZFNzdmQKNjlCV1Q5NVlEN1RsU1RWdGhnQjlwbzFpVG5FMnd6VXpQSmdsNWZTMGpqdWYwZFJVMHloSk5wZkdrYUxwTmJWVApBSHZDU3lBUTJrSXlvZTFkeWdQS2ZEdU9aQTZnTWF5K1p3ajNIZlB4SmJxUXNmWkJWMUkwblRXZGxsYlBhdGQ3CmRjdERNbit4NkErcUJNZnhxT1RiZm9FM3YwTmhNQXIrdU5YaTAvd3F2UzBEVUxVR1VxWW1lT3drL3pYSGRNdS8Ka0NtTGxFc0NiZ2Mva3NZbEtvM0lUMnV3RlNtQ0xqa3U5VnVkV0lFa2hwWkRrZUlDUlNRcm1ISGpFVjRMNzZ6bgpsWGd6c3ZaM2p5VzRQSkwzMnE2NWhpazVJTkpNNDdmL0pxVnBJZ2gxeGQxWnJReTZISE5rWU9DVzI3R3h6cUg4CmI2Y3pqZGQwZ2dMQXlFRDM0UmFJazNBYkVldUtyeTR2ZGVxVHFzeU9HdnFCRUdYYlVibEtPMENhTXFPMi8zSW8KWlZTRTNOajlYSURxQUVpSnplc3g1b3FaWW4rQzlOQm1JZXNMU1JBdTZidk5KZmdleVJkVWgvcnliYzZmMWJLUgpFSURYQnF2QkF4OWpNUXE5aVlERGhvSVhyNnA4anVCdFdPS1YxVVZObVV3WGtabGh3WVhkMGhhVSt4VnR1SzhKCm04ZTU4aUIwV3MrczBrV3ZMNVY4aUlSYVdEZ3pmdHYwNFgvL0VzUjAzcFd6OGhTLzV6dEJyMU1wTHl2NHFjNGYKOTdZbnQ2ZnN4dnF3eGx0cmtXSUNsdmRobk9Sa0xDa0ZnY3Fna3JydDZoSGFzS2ZEcGprSjExQXVOS1drbU9wUwppUHEyajArTmFyVmFSYmtrVmcvVkhXSlB5UXNwNVFJREFRQUJBb0lDQUc1QlZVUnZLem00WHlMRkNTSThCZDNIClhLa1FTbG5GRkpPK2lydHRrYzJVQzZpRmxhanJIbGoyRk14ZEFLUC9uNjVZZTVDekpZTzFsSWxyaWpBVWV1L2MKTlRDMGtDY0FSeWY4ZWFMQ0lkWlJuYmhzTm93ZXJFZTE2U2dpRVRFK3BoWkorYThBZ1o2eUx5R3ZSNnhWMDJYdwp6QUtsU0tmZDZEaDRTMDQ4clc3YXBIZnRGVWZ1N0FuaDNnNS8zN0hOZ0NySEpiODJtclZPSDdGOVhmdjJ5N2tvClpXa3pWRm1iNGMrenBxV0JOMDRSeFo2N0hNODltdlNQemlCd3VseVRkUXpGNXB3VkU5WEJ1Z3JOVHFDU3dZOXUKZU9qbHJmUFlxd09UbFBpMmNCQWRlc0daYmxVT0d6c0dwN2VEWk9oaVhLZDQyWDZ3bXNDeEZlbDNPS29rMXkyQQoxN3h0Sm5zbnMwdFp1K0FDU1dtdUkrbThkblZLcXgyVjJBWTRPd3lPclJ4dldxTmF2WnlmVnFnQzdQcnIrNkhnClcyd3JWN0EyWk1iVExsNkF0TU80NzJNTStOdnNtNHhvRDcraVMrNVNUMVJSYVR6MWJzZ21BWUxMWFF1ZnFhQzIKYStVLzM0elRGaEthUWFaekZUdnA4Ujl0R0VMSTlqY1pqWGxyT1QweFlkZ3lDTUMzRXNGSXU5NkJuREFGWndLdQpmTzN1dDJxY1NBOXB1eUpKUVRNZEtxSlpLSUlWVnJHVE54aWJXalB4dU5ET2VFVXh3cll5UnJOemdlUVdQWEpjCnVrcnZvdFpUSm5uNktEODdZYUdrUm9WOFlia2lWY3JuT3k1N1dlOEZhRnZtaWsvZ204Qnk0TllwZGN6WDUvOVEKblZPU2RxSDh1bTVNR0dRT2dSeUJBb0lCQVFEVnpuWTdMWVVvUXNWR2RVVEVsQVFLUjNhSUxPK3k0aFdWblFDZQphZllOT0NyVmFpVTVZMG9uc3d2c0xrUTNEV01ZTjluYkxCdjkySUZhRTJISmN6bUtlRVFvTzZSVmVCeG05ckxlCk5KVURTWkl1SkJqTEx4WUJzV1c3b2s0WUdnZFZpbEVCeUZuN3c0V2tIcWlJektPSmx2Rkd1bVlEeFdJSlA0MDEKV2t2bzFqUWlOYXZpUkNHeTE0OGVMdUJTUDhsSWlBZm8xQktsVFVob2FmbFZXUCtGTlBML0R6R3VidDNuK21OLwpMV1M5YUVaaExDNkFJT3lhbVBlRVN4eW1uZEpWblNKVHZKK1NabWR5UVNuSXZML0o5NUNJUDdleUlUYnJTWlZxCjNpV1JLaURVUTYwUTgyNHdZaDdaYVlBYjVuZFd0QmVFT0dTRHhZQjk2c2MvTldEcEFvSUJBUUROTzJ6THJnaFEKdWY4R2ZpcWdXK1lESTczUEQzY2U4dGw3MVdCSVkvb0FXOWQra2tHamtvZ2RqLzI4WVZEOVRKK295dVRSVzRQQwpTMDNPb0ZyK3U4S3FGcXRXVko5TUVHY0ZTSGdCZms1TG5JR0R0UCt1dnUvTHZGTVEyekJkRlo4bSs5S2tOdDd2ClA3aWx6a2pwN0U2TFVTYjNqMzc1YzYxaFBMMVJhOERpeGM1d2kvOXIwbndyanVjRE9uSklOOThYMGd6YXBINFUKMkMwM1NRUmRNQno3VjVSaUdHbTNRaGRCWHBhb2tiWVZneklQbU5zd3pnZHdLN2tFRDdBeHM0RFFwMjAwMEZ6Qgp0M3pzQVBQRXMyRm5CcXVXUlczUjlhd0lyaVhNL0x3dXRvaFhhK0VrbU43eUcxK0VRZEVOQzFRcFdrTElERXY3Cm16bC94b2FnVEFPZEFvSUJBRHcxNEdYWjg5M2FyK09mc3JZSldQbnNGaDFUU2sxK0RjWU1hTmd6enU3NkdsWHYKaG53YTBnOU1CTmVHVC8rUTdZOHNhMVdsbmx4bVZFY2huakExR3NjOEJ6V3RWaUlicVNQMTVYbGVKWGkvaDBNbgpOelJCRmxsenM4cWJjcEtuQWRtOUVnTUdnUkM5aHkwbzFSMXhRN3pEblQ3bHowVFFtVU14ZW5yRDZ3eXZCZzk4ClBlT0NmRnI1Q1h6ZWhwMmpDUFE4R3I0ZXV3R0NPaG50ZmlIaTVsS0ZEc2wxWmZCUm1IeHpyd0ZwcnkwSDZJb0UKL3pObUVqdVhTRjBoS2ZoaUNaSENwcUFlUm5IY0ZOWEFOQndyeTNiOUdON0YwdDEvTFJBbHNNWmZ1UVNnY0k1VwpZSzZkWHpLUTcwOGF2dEVjbmc3MHVJcXJ0dUxGQStKeDg1cUJWY0VDZ2dFQkFMbmdSRjBFdGd0SEthN2J2Z2VXCnMrL01BekR4dE5XVzVWcSttb0YxNnd0QUl5QkRucWRqSTF5QytUQVFnNldtTEVSWDNuMnZBTnFNRVdBKzQ2c3EKcXRnWngveGNrQm40RVJZNzJGU2g3SStXbzhhQnU4Q3N0Y28wT3BkZHJhUGczVkFWYTJYSFBJbzdrQ1M3ZkZaQQo2N0pLUVp5ZG5rYVhla1JER3NRUGI1Ynp5RkV1dXBzUSs5MEhoRHJzU0cwWURUb3B4L0tPWUpMSVo4dFdtbGs4CkprT053cHBGdWhsOEJrdnlPMGxaRHl6VXNoWm1QcjhwR3B1QlBnUnJvUXlpb3R4WGh4VDZVY2d4UXpjTWRidUYKSzRhQUNCQUZ1YjBiWUVCTVdYZ2F2dVVmOU1RWXRNVE1uNzl2QTBkcHhNaW5wZ1g5OWRYeExUQW9HaCtiMG5xRAozaEVDZ2dFQkFKeERNcGJhVTREY2RxUG40djRZNVR5amhGWVVxNS9ZQW9vNFgwOG4ySHFPZzc4Mmhydm9jR2pyCm9IZXh6NGVjWjY1OHFMY3R5VGoyV0h1dysrNkVTcTVZUDg1bFdNc2VDTXhteHBIalJhWnc1dW1YTEpROWYvUC8KUkgzbkpqUUw0OEhCWkRKYi9MUFJSTXZrc3ZQY20vRHQ4ek1qSVVJT3lhcHZOK2pKSWxuRTArRXZyY2krV1R2cAp2WE5LRVVDU21KbDJpUk41Umg1RVlYUGJpc0tEZlRpU0N4OHFmV0ZDZ2djejR3S3NNandwR0JuVllWeEVtcldFCkMxVWRuZmR2SHUrR3FYYXNVY3RVbWt4M2pLODVPNkswZUFvWFRzdGxqWkVwRzM2QUxySVdIOVExRWdtOWk1RUkKbVhuSlFNTHo1WU53SkE2Z043YWdLYW5YeWdGOXdPbz0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQo=
kind: Secret
metadata:
  name: cluster-ca-cert
  namespace: ibm-mongodb-operator
type: kubernetes.io/tls
`

const mongoCertYaml = `
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZiVENDQTFXZ0F3SUJBZ0lVUjkwRTZSdlVmb29sc3N6NTZ2V3dEUnpuL1pRd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1lqRUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdNQ2tOaGJHbG1iM0p1YVdFeEZEQVNCZ05WQkFjTQpDMHh2Y3lCQmJtZGxiR1Z6TVJFd0R3WURWUVFLREFoU2IyOTBJRWx1WXpFVk1CTUdBMVVFQXd3TWJXOXVaMjlrCllpNXZibXg1TUI0WERUSXdNREl4TnpBNU5UQXlObG9YRFRJeU1ESXhOakE1TlRBeU5sb3daekVMTUFrR0ExVUUKQmhNQ1ZWTXhFekFSQmdOVkJBZ01Da05oYkdsbWIzSnVhV0V4RkRBU0JnTlZCQWNNQzB4dmN5QkJibWRsYkdWegpNUk13RVFZRFZRUUtEQXBFYjIxaGFXNGdTVzVqTVJnd0ZnWURWUVFEREE5dGIyNW5iMlJpTFhObGNuWnBZMlV3CmdnSWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUNEd0F3Z2dJS0FvSUNBUURLWHdPb0M0ZG1aeWJEdzFRUWFxTnEKVHUxbWtGcE9UK0oxYU9JM1F4MHc3dkR1bkpWbjd0aVVGS1BMRVhiYlRHcUthZk5hVEZ4ZU1mQmZUTmUyWk1JQQpCTjVGdFFPZ2pUblBOS052T0dwWUNMbGFibG9lZVpoUTEvQkhTNHQzK01SbDYzSWtnUkFtMUVnOFhWMi9xMnFLCklLQVNNTmZvdVV1ZkU3NlV5ZVVhcGdhQkQ4MWlLWDRTUW42T0FFb2taL0hLaTh1b1RDOFRpMFFscENHbDlSTUwKVjFvMFRzVGViQ29qVHNMZWUrZFN4VjMveUJmM1RzaXltUmxndUpIeiswazVhYjN0QXFTb0lFVW92Vnc4bnljRQpGamMvLzhGanF2VVUzb0dEZXo0dkt3UUhYRi8rTWt6dHkvSGUxL1o0bFcwNnBXbi9vNWpsYXhITHVyS2ZKeG12CjJKTUg5ZTVueFl6RkcyMHJXT3l6elhQdWtZanRUZjdrdFZjVmE0U1loWFNRc3hBVCtsUnNQYi8wcUJuaWhIdkYKa1AwRkpqZFdTaTNTTmx3MWtFSTFaaDRHZkEyejJpMVlIZjJlcXhFNm0reTBEYkcwVGIzL2JDbnBJWFpBdlBmaQp3SFlPbUJWdXpnUVIwd2VmU1I4MFRYWktqK1VReWFEM2h0MXg4ZE9ySnhvdVhmSFZwZ05KanVFaEl1SSszQXM4CjlRZlZISCtNZjh1Q0ptUEdwK1dMaHpGTW1jMHNJYmJGSHJ0V2liMVh3R0VSVjM3bFhPNkpJaDdHNHNSQWp2YWkKVUpzWFFFNysxMldraUZUc241MWJ0OFJ0dXZhSklKWDltS3JjS3lhQmNBb1lnTHNlMWpjV3l2YzUyd0pTck4yZwpXS2prVGs0RlJFYlBvNk9MR09aN3N3SURBUUFCb3hZd0ZEQVNCZ05WSFJFRUN6QUpnZ2R0YjI1bmIyUmlNQTBHCkNTcUdTSWIzRFFFQkN3VUFBNElDQVFDV3dNU2pHZHU1cnB2MXB4SVF1WVJ3Vzc1c3Y0TTYwSmxIUnFiU2hwT2MKaFdaR3J3cmtTWStqREVpckNaaGxqL1Mrd2xmSXE3ZUtxV3hwbnIxOVpONENaWk9uL2lqL3pCeStSdWxkVGxzbgp2VENYSHlDVmE3S1ZjaUF3a3grQVg0V0p5ZmJFTHkycEtvaXhaUmg4SHBYZ0s5NlZCUi9RbldVNHFnSWZKUUZlCnRzdnp4aE5ya3BCajVTa3JrcHZjUEFzMHNqd0tQRVZvRjkyOWtJUjNCSFFlK25lOUEvSkhKaEdyOU5mR2tua2UKNWFrb0UzQXRQYXJlbCs1R3FkZzNIRkUrbnNPcWdwSWdjRlB3N1k1dk5UWThrTWZFcmpTT3QrN0JMK0M4Z3EySApvb1Vkb1g0YktSTTNmVmppalV3TVZJY3pRYmowNWg0VHhabVR1T3pLdk5NTXVQS2JjdzUrMzR4Yk9SbUtmWUZxCmt2MDVJZ1pHdkcvYXRUTXBmNHlwSTV3Y3M1VkxJL0M5OXlndUdiSHJjTElWOTJSUzdmYndFc1IxVGVoQmtvUkoKNnkyek1yTGRVSFVOdFE1VlZKTHhUeWN3bGlFUXlzT0l5ZWcxWWFxN291eHhMeGgvRnQrd1VkdGRienRaVE14NApyNTdtNEtSdnZteDlPYkR3VDZYdHVrRzVZUTFHcVljUzZGbkdyRU04bVcwR3YyRmFJaytIMHppWTA3NXVEY3ZQCi9nM1BhUGU4UUxpVUpwL2ZXKzBwVXhyMnFTTGlUWWkrMjdaWFRoNjY2VW15NWFmc2U3K01DMDBSUGdueWk0MEQKUEFmamxnWWM3Q3F3QWh2cU9XN0JhcmI3V3drSGJMbUF2ZWZsOGVlZGMzNnc3T0lWNFlTOWg5S3pscXNtNVdJUQphUT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS0FJQkFBS0NBZ0VBeWw4RHFBdUhabWNtdzhOVUVHcWphazd0WnBCYVRrL2lkV2ppTjBNZE1PN3c3cHlWClorN1lsQlNqeXhGMjIweHFpbW56V2t4Y1hqSHdYMHpYdG1UQ0FBVGVSYlVEb0kwNXp6U2piemhxV0FpNVdtNWEKSG5tWVVOZndSMHVMZC9qRVpldHlKSUVRSnRSSVBGMWR2NnRxaWlDZ0VqRFg2TGxMbnhPK2xNbmxHcVlHZ1EvTgpZaWwrRWtKK2pnQktKR2Z4eW92THFFd3ZFNHRFSmFRaHBmVVRDMWRhTkU3RTNtd3FJMDdDM252blVzVmQvOGdYCjkwN0lzcGtaWUxpUjgvdEpPV205N1FLa3FDQkZLTDFjUEo4bkJCWTNQLy9CWTZyMUZONkJnM3MrTHlzRUIxeGYKL2pKTTdjdngzdGYyZUpWdE9xVnAvNk9ZNVdzUnk3cXlueWNacjlpVEIvWHVaOFdNeFJ0dEsxanNzODF6N3BHSQo3VTMrNUxWWEZXdUVtSVYwa0xNUUUvcFViRDIvOUtnWjRvUjd4WkQ5QlNZM1Zrb3QwalpjTlpCQ05XWWVCbndOCnM5b3RXQjM5bnFzUk9wdnN0QTJ4dEUyOS8yd3A2U0YyUUx6MzRzQjJEcGdWYnM0RUVkTUhuMGtmTkUxMlNvL2wKRU1tZzk0YmRjZkhUcXljYUxsM3gxYVlEU1k3aElTTGlQdHdMUFBVSDFSeC9qSC9MZ2laanhxZmxpNGN4VEpuTgpMQ0cyeFI2N1ZvbTlWOEJoRVZkKzVWenVpU0lleHVMRVFJNzJvbENiRjBCTy90ZGxwSWhVN0orZFc3ZkViYnIyCmlTQ1YvWmlxM0NzbWdYQUtHSUM3SHRZM0ZzcjNPZHNDVXF6ZG9GaW81RTVPQlVSR3o2T2ppeGptZTdNQ0F3RUEKQVFLQ0FnQWtBdEJWd09keE0zM1ViQmV6YkNaMExtTlVVdStlNjl3eVpGMk0wK2FINUp6KytPSWxRbjFMckhpUgpGQ1NBVlpMSDJwNnhQTkZhK2F2NmFXUWhVc0NxM0RMcFdKS3lxUzdXVGxtZTJ2MGhlVHZ5ZVp5VHU3TjgvMUFFCmY4N3JwRnJlZ0EwcHJjWEFBeHB2azNXeE84R1Ruc2FkTmcvVm05TjNGVDVlbjZhakhWUWU5ejdtN3RjK1RKTFUKbGZ5YmlkdWUzVTE2UDBSSlNBanlZY2lURFk4Ny8ybFAwWXg2ditpbnE4WkZiT3IyOGFRT2RmNjl4VWsxYnNUegpUeVM1czhlTjdlRWNJZEpIRUtiOTN1Umc3VGsySXZYbDc0N3NPMm10TXdMODhKdGFMVjlrSiszMC8rSnNsbFFPCkFZUWNaUXF3MnVxSDBRYk9IRVZvYVdxTG81dVFQQVpzL29ENlViOGZyRFNKZnlvWUdxWEdkZTBXMmlGcXEwZ2UKeHlQZVN6YUlxenprdnZZa1RJbUNicDhUK2hSeVY4RFBndU44VUltQnlEdUNNYUZSVVBYVGF0bytUcTN4T1BrQQo0b0ZoTlJQNkJ3VGlDS3NQenpJQ2MxdlczRHFxWFkwWTRrajR6ZXdzcHJLakZGcEFvQnAwQkVab1NXUk5hVEZ5CllDSzFKQTdlMnVYbUNuNEFZTXM3bEI4dmJVNFVjeXAzOVVPNUlEangvMEdvYnhZSXZ6clg2c3NLRTZwb2ZsaUoKTkZWOVNKVGZITU1NN1RqaVdGQnVSbHNmQ3NiMVRpT1hXcXl2SElWM3hremUyWUhmU1ZRYTk3WDE5TzFzUitTUgpsaytuY0l2Mk1jdWdDK2JHTENFQ3dqV0VGaXZ6R2lsK3p3blV1VGpVdnplTE82a2R3UUtDQVFFQTUvNWlkZWM5CmllcXdlNGdYQnF1emFURzdSY2pjcDl6Y3dDK3d5VzYrSm9tdzFnNXdSWDNWVXJYVzN6dDFFWWVYR1RJVy80WHQKT0dxTjRSaERNaHRrS1hFRWt6eHhDbmh1MTZ4K051QTllb1JQdUdzdE1xbHhzV25NK2RqUVJKdnNnVlhSc2VNdQppdkc2cHpNd0ZQZFR5ZWlYcGZEMVBBOXdLQVFKSmowdnVmN0daUkw2ckppajc0Rk91eHA1VGgxbzZ6bjJwaE5MClBTZFBVWXNjMURhbU1pcGxYMm9HS20rQTJlUlhYK29RMmdia2lMS0QzSWRCbjBUU1AzZEc2TU5Cb2J6NW1aVTQKMHBTTU9iYi9CNytQMU11aGcrYkJsUlc4LytEbXNCWXJSd1Q4UVFTVHBNcUpCbUlBa1JnNFR2ck1TYmtjQ1o3cwpWNjNBK21GK3hiSllnd0tDQVFFQTMwL3FYU3RkK3AvbDdSMjRuU2x4Qi9uL3dvLzVUejRSVm0wY05YK1ZvMjZzCnVKT0h5WDlJT29JMjVQcExBT0pDZThNNExiMVRZWE44MTZFN1dwc3ZJOWhyWk9WQll0MXZwSHdGL2JETUt4V2MKSi8zQnhUbXZjUXB5TE5WSGpRRVJybTY2Tm5SanJCeFdmbGIydGgvRy80NThpNjhVVURrQWpmdUlHNE1NZ256Wgo4enYvWEl5YjV3aHlPanhqZ2ZSSmxYdjI5MnJ0NDZJV0hVQWJOblNOVW4xVVAzV0dVWEh2NmdaM1RydVQzRmJCCjNoSzRnd3F5WTZtOHBObStCK0t1Zjk4WHd5eVYxSmtuMnB3OENrTXlFU3ZsQmtQRXEvV0pEV0NhT3Z0UDl5aXoKR3hvSU1qSS9wT2NWZXhTQnowbVQ5UVZHQ2hwYnVNMUNTbFV1R3FZSkVRS0NBUUE5T0FXbmMwUHI5d0JuT0x4Swp5Rmhwcy9QbE1HSDU4ZkJXenI3cUNNMG93a0RsMjUySTJQSElCN0FSN0ZDeU5ZT0w1SW5wRitCSGVPYkR0WEZWCjhhQjJ4eG9iK0dFa0VDKy92Z2I0V0NnaEFuVS9CeGxBT3pLRFRKWUlnRXhGTHBnMGNQOEs0QlpTR0FQWFIweXkKMjZsQ3FKd0w4QS9tcjNRN093Vm5EOUplVkhycUJSNGRHWko1Q3poSmEyMERUZ04zdnkzMUdUWkxodW9KYkpwSQo3YnJobGdwMktUWkRVSFZDQ2wxOE0vb0tick16MTFld2hBaXZES3dtajBVbyt3MkFycXQyK2NlcTJnUSszcWxoCjFBMFNiRUhNMnNIT281UGlPZWptSXBOOUJEWEV1bjV6aC9hc3RvUEx4Z1psNFF5emo2TjBibm1Ua1loUkNoVTMKK2g0ZEFvSUJBQ1VWUFVMNWg2S3QyTjIyV01qb2I5ZTJRUzJMQVFpU3N3aGFHQndlTXJnd0VjaVkzeXlyMFUrOQovZVdxVnJndjJvQjQyNlJrMHlyVXBiK2RDNkV4TWZQTzVZNmNyMjMrZmFLZjRkTE9BQ21MYmlJSjlwcU15TUNKCnpvbjVaT2RhYlJnOVZQamovUVZBczNCSmVyQ2x1RU1KNDA3QzVTbXBQWmxXVXJUVzMwWHYrN1Z6bWlWQlNFWm0KVmFtc0M0NHlCZUluOHN3RldybTVXZGpEbzRFNGU1dGVLcFpiS3RIdGpMeWRGRVRqeTFzRW9TOENodGRqK0ZtcQpmeVFVOElTWXRRZVJBWDRzc2pqYXNnNlFjVHYzQ3FKbFdxUGVyeE1yTS9ZZnU1emR6TnFyVElyTW1OM1ZFRktPCitUYzJJWlJOa0o5WW45ZmZwcW1hbEU2SnRKMUNRekVDZ2dFQkFPV200UkFTbmRZWjNXQ3VQWkQ2b0VzckI3ZTMKWmFDT3UrK00wMHY4dWJqbzVDbUd0S2pDSW53TnVnSmx3ZVBwUDBlQTF1SHQ2YTdjSFh5MklkZlJscEgxZ05DZgpVYXRGWVJlTDBKVk1kOUlERGhTVVAxSmtTdG5lN0JRUVRnRXA4SzhSbnBoRDBIZmQwRlc0NFRlUHlVdDlPWmY0ClpTamxOOEtZNWdHcTJ1WnM4dHZsWWUwS0N0VEZQc2NuMW1RWmdqYmZoSjcrNnA3aXVDek13bEE2Z3VlWStOT1kKdTJtNWNoaFd1K2xlMjhDUWVZTEM2a3ZvOW5QVjBJNDBEanlEd0dOTEYyNHRsMVNYRStLcXJXZzVRNSt2MzhVOApGazN4RThoeHV1TVRBODlYQTBEMDluWlJ6N0dXUmhuSk5MZmFhOFhZSHh5QXVmT2xTT0pRYW13Z0dVQT0KLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K
kind: Secret
metadata:
  name: icp-mongodb-client-cert
  namespace: ibm-mongodb-operator
type: kubernetes.io/tls
`
