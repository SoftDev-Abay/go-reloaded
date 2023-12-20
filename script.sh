i=0
while read -r line ; 
do 
    i=$((i + 1))
    if [ $i -gt 11 ] ; then
        break
    fi
    
    echo $line >  testdata/output_tests/test_case_${i}_output.txt ;
done < result.txt
