# [ Programmers - Prime Question ] - [ 내가 생각한 핵심 ]

>
* 완전 탐색 구분( 찾기 문제이면 탐색 문제 인 것을 염두해야함)
>
* 주어진 값 중에 가장 큰 값을 구해야 내림 차순으로 값을 놓는다
  - 예를 들어 1234이면 4321 []int형으로 만들어 놓아야 한다 
>
* 가장 큰 값의 크기 만큼 []int 배열을 만들어서 소수인지 아닌지 구분할 수 있도록 만든다
  - 예를 들어 1234이면 4321이 가장 큰 수 이고 4321 길이 만큼은 int형 배열은 만든다
>
* 그리고 범위의 값들 중에 소수인지 아닌지 구분한다
  - 소수인지 아닌지 구분하는 로직을 넣는 것이 아니라 이 값이 이미 소수가 아니라는 것을 전제하에 로직을 만들어야한다
  - 예를들어 2라는 숫자가 소수인지 아닌지를 구분하는 로직이들어가는 것이 아니라 2는 이미 소수이므로 표시해주는 로직이 들어가야한다
    + <이 내용을 잘 이해해야함>
> 
* 해당 범위의 모든 소수와 만들어질 수 있는 수열이 일치하면 소수이므로 갯 수를 추가해준다
  - 해당 수가 선택 되었는지 안되었는지 구분하는 로직이 들어가야 한다 예를 들어 4321 이라면 4가 먼저 이미 선택된 값인지 체크한다
  - 값이 체크가 되어 있으면 다음 수를 진행하고 아니라면 선택되었다고 picked 배열에 표시해 준다
  - 그 값이 소수인지 아닌지 체크하고 소수이면 결과 값에 추가해준다
  - 그리고 이어서 재귀함수로 또 로직을 실행하고 방금 선택한 수를 인자로 전달해준다 
  - 그리고 결과 값을 리턴해주고 다시 해당 위치 picked값은 false로 다시 만들어 준다
  - 결과 값을 리턴 해준다
---
* 생각을 못한 이유 
  - 우선 모든 수열을 비교할지 생각을 하지 못했다 
  - 수열이 소수인지 아닌지 구분하는 로직을 생각하지 못했다(사전 지식이 굉장히 중요했던 부분)
  - 그 만들 수 있는 값을 어떻게 모두 만들지 생각을 못했다
    + 예를 들어  4321 이면 43,42,41,413 등등 이 수를 어떻게 만들어서 비교할지 몰랐다